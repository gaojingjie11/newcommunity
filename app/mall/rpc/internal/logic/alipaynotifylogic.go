package logic

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"smartcommunity-microservices/app/mall/rpc/internal/consts"
	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AlipayNotifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlipayNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlipayNotifyLogic {
	return &AlipayNotifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlipayNotifyLogic) AlipayNotify(in *mall.AlipayNotifyReq) (*mall.AlipayNotifyResp, error) {
	if l.svcCtx.AlipaySvc == nil {
		return &mall.AlipayNotifyResp{Success: false}, errors.New("支付宝服务未启用")
	}

	// 1. Verify sign
	if err := l.svcCtx.AlipaySvc.VerifyNotify(in.Params); err != nil {
		l.Errorf("支付宝异步通知签名验证失败: %v", err)
		return &mall.AlipayNotifyResp{Success: false}, err
	}

	// 2. Check trade status
	tradeStatus := in.Params["trade_status"]
	if tradeStatus != "TRADE_SUCCESS" && tradeStatus != "TRADE_FINISHED" {
		l.Infof("支付宝异步通知: 交易未成功, 当前状态为: %s", tradeStatus)
		return &mall.AlipayNotifyResp{Success: true}, nil // Return true to avoid Alipay retrying for non-success statuses
	}

	outTradeNo := in.Params["out_trade_no"] // This is our OrderNo
	tradeNo := in.Params["trade_no"]        // Alipay transaction ID

	// 2.1 Handle Recharge Callback if prefix is RECH_
	if strings.HasPrefix(outTradeNo, "RECH_") {
		var userID int64
		var amount int64
		err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			paymentRepo := l.svcCtx.PaymentRepo.WithTx(tx)
			record, err := paymentRepo.FindByOrderNo(outTradeNo)
			if err != nil {
				return fmt.Errorf("未找到对应充值记录: %s", outTradeNo)
			}
			userID = record.UserID
			amount = record.Amount

			if record.Status == consts.PaymentStatusSuccess {
				l.Infof("充值记录 %s 已经是成功状态，忽略此次回调", outTradeNo)
				return nil
			}

			// Update payment record status
			affected, err := paymentRepo.UpdateStatus(tx, record.ID, consts.PaymentStatusInit, consts.PaymentStatusSuccess, "")
			if err != nil {
				return err
			}
			if affected == 0 {
				r, err := paymentRepo.FindByOrderNo(outTradeNo)
				if err == nil && r.Status == consts.PaymentStatusSuccess {
					return nil
				}
				return errors.New("更新充值支付记录状态失败")
			}

			// Call WalletSvc.Recharge to credit the user balance
			walletSvc := l.svcCtx.WalletSvc
			err = walletSvc.RechargeTx(tx, record.UserID, record.Amount, outTradeNo)
			if err != nil {
				return fmt.Errorf("增加用户钱包余额失败: %w", err)
			}

			return nil
		})

		if err != nil {
			l.Errorf("处理支付宝充值回调失败, orderNo=%s, err=%v", outTradeNo, err)
			return &mall.AlipayNotifyResp{Success: false}, err
		}

		l.Infof("成功处理支付宝充值回调, orderNo=%s, tradeNo=%s", outTradeNo, tradeNo)
		if l.svcCtx.EventBus != nil {
			l.svcCtx.EventBus.PublishWalletRecharged(userID, amount, outTradeNo)
		}
		return &mall.AlipayNotifyResp{Success: true}, nil
	}

	// 3. Process GORM transaction
	var order model.Order
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		paymentRepo := l.svcCtx.PaymentRepo.WithTx(tx)

		// 3.1 Lock order row
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Items").Where("order_no = ?", outTradeNo).First(&order).Error
		if err != nil {
			return fmt.Errorf("未找到对应订单: %s", outTradeNo)
		}

		if order.Status == consts.OrderStatusPaid {
			l.Infof("订单 %s 已经是支付状态，忽略此次回调", outTradeNo)
			return nil
		}

		if order.Status != consts.OrderStatusPendingPayment {
			return fmt.Errorf("订单 %s 状态异常，无法完成支付, 状态: %d", outTradeNo, order.Status)
		}

		// 3.2 Lock and update PaymentRecord
		idempotencyKey := fmt.Sprintf("pay:alipay:%d", order.ID)
		record, err := paymentRepo.FindByIdempotencyKey(idempotencyKey)
		if err != nil {
			// If for some reason payment record wasn't created initially, create it now
			record = &model.PaymentRecord{
				OrderID:        order.ID,
				OrderNo:        order.OrderNo,
				UserID:         order.UserID,
				Amount:         order.TotalAmount,
				PaymentMethod:  "alipay",
				Status:         consts.PaymentStatusInit,
				IdempotencyKey: idempotencyKey,
			}
			if err := paymentRepo.Create(record); err != nil {
				return err
			}
		}

		if record.Status == consts.PaymentStatusSuccess {
			return nil // Already updated
		}

		// Update payment record to success
		affected, err := paymentRepo.UpdateStatus(tx, record.ID, consts.PaymentStatusInit, consts.PaymentStatusSuccess, "")
		if err != nil {
			return err
		}
		if affected == 0 {
			r, err := paymentRepo.FindByIdempotencyKey(idempotencyKey)
			if err == nil && r.Status == consts.PaymentStatusSuccess {
				return nil
			}
			return errors.New("更新支付状态记录失败")
		}

		// 3.3 Mark order as paid
		now := time.Now()
		result := tx.Model(&model.Order{}).
			Where("id = ? AND status = ?", order.ID, consts.OrderStatusPendingPayment).
			Updates(map[string]interface{}{
				"status":       consts.OrderStatusPaid,
				"used_balance": order.TotalAmount,
				"paid_at":      &now,
				"version":      gorm.Expr("version + 1"),
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("更新订单状态失败")
		}
		order.Status = consts.OrderStatusPaid
		order.PaidAt = &now

		// 3.4 Record wallet transaction (for bookkeeping, type: OrderPay, Alipay method)
		walletRepo := l.svcCtx.WalletRepo.WithTx(tx)
		wallet, err := walletRepo.FindByUserIDForUpdate(tx, order.UserID)
		if err == nil {
			// Record transaction but don't deduct balance
			_ = walletRepo.CreateTransaction(tx, &model.WalletTransaction{
				UserID:         order.UserID,
				Type:           consts.WalletTxTypeOrderPay,
				Amount:         -order.TotalAmount,
				BalanceBefore:  wallet.Balance,
				BalanceAfter:   wallet.Balance, // unchanged because it's paid via alipay
				RelatedID:      order.ID,
				BizType:        consts.BizTypeOrderPay,
				BizID:          order.OrderNo,
				IdempotencyKey: strPtr(fmt.Sprintf("pay:alipay:tx:%s", order.OrderNo)),
				Remark:         "支付宝支付",
			})
		}

		return nil
	})

	if err != nil {
		l.Errorf("处理支付宝支付回调失败, orderNo=%s, err=%v", outTradeNo, err)
		return &mall.AlipayNotifyResp{Success: false}, err
	}

	// 4. Publish Event
	if order.Status == consts.OrderStatusPaid {
		if l.svcCtx.EventBus != nil {
			l.svcCtx.EventBus.PublishOrderPaid(&order)
			l.Infof("成功处理支付宝回调并广播 order.paid 事件, orderNo=%s", outTradeNo)
		}
	}

	return &mall.AlipayNotifyResp{Success: true}, nil
}

func strPtr(s string) *string {
	return &s
}
