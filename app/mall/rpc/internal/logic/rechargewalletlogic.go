package logic

import (
	"context"
	"fmt"
	"time"

	"smartcommunity-microservices/app/mall/rpc/internal/consts"
	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type RechargeWalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRechargeWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RechargeWalletLogic {
	return &RechargeWalletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RechargeWalletLogic) RechargeWallet(in *mall.RechargeWalletReq) (*mall.RechargeWalletResp, error) {
	if in.Amount <= 0 {
		return &mall.RechargeWalletResp{Code: 400, Message: "充值金额必须大于0"}, nil
	}

	if in.PayType == "mock" {
		idempotencyKey := fmt.Sprintf("recharge:mock:%d:%d", in.UserId, time.Now().UnixNano())
		err := l.svcCtx.WalletSvc.Recharge(in.UserId, in.Amount, idempotencyKey)
		if err != nil {
			return &mall.RechargeWalletResp{Code: 500, Message: err.Error()}, nil
		}
		if l.svcCtx.EventBus != nil {
			l.svcCtx.EventBus.PublishWalletRecharged(in.UserId, in.Amount, idempotencyKey)
		}
		return &mall.RechargeWalletResp{Code: 0, Message: "success"}, nil
	}

	// 支付宝充值流程
	if l.svcCtx.AlipaySvc == nil {
		return &mall.RechargeWalletResp{Code: 500, Message: "支付宝服务未启用"}, nil
	}

	orderNo := fmt.Sprintf("RECH_%d_%d", in.UserId, time.Now().UnixNano())
	idempotencyKey := fmt.Sprintf("pay:alipay:recharge:%s", orderNo)

	record := &model.PaymentRecord{
		OrderID:        0, // 充值没有商城订单，设为0
		OrderNo:        orderNo,
		UserID:         in.UserId,
		Amount:         in.Amount,
		PaymentMethod:  "alipay",
		Status:         consts.PaymentStatusInit,
		IdempotencyKey: idempotencyKey,
	}

	if err := l.svcCtx.PaymentRepo.Create(record); err != nil {
		l.Errorf("创建充值支付记录失败, user_id=%d, err=%v", in.UserId, err)
		return &mall.RechargeWalletResp{Code: 500, Message: "生成支付记录失败: " + err.Error()}, nil
	}

	payURL, err := l.svcCtx.AlipaySvc.GetPaymentURL(orderNo, in.Amount, in.ReturnUrl)
	if err != nil {
		l.Errorf("获取支付宝支付链接失败, orderNo=%s, err=%v", orderNo, err)
		return &mall.RechargeWalletResp{Code: 500, Message: "生成支付宝链接失败: " + err.Error()}, nil
	}

	return &mall.RechargeWalletResp{
		Code:    0,
		Message: "success",
		PayUrl:  payURL,
	}, nil
}
