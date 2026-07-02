package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"smartcommunity-microservices/app/mall/rpc/internal/consts"
	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentService struct {
	db               *gorm.DB
	orderRepo        *repository.OrderRepo
	storeProductRepo *repository.StoreProductRepo
	walletRepo       *repository.WalletRepo
	paymentRepo      *repository.PaymentRepo
	userRepo         *repository.UserRepo
	productRepo      *repository.ProductRepo
	eventBus         *EventBus
	log              *slog.Logger
}

func NewPaymentService(
	db *gorm.DB,
	orderRepo *repository.OrderRepo,
	storeProductRepo *repository.StoreProductRepo,
	walletRepo *repository.WalletRepo,
	paymentRepo *repository.PaymentRepo,
	userRepo *repository.UserRepo,
	productRepo *repository.ProductRepo,
	eventBus *EventBus,
	log *slog.Logger,
) *PaymentService {
	return &PaymentService{
		db:               db,
		orderRepo:        orderRepo,
		storeProductRepo: storeProductRepo,
		walletRepo:       walletRepo,
		paymentRepo:      paymentRepo,
		userRepo:         userRepo,
		productRepo:      productRepo,
		eventBus:         eventBus,
		log:              log,
	}
}

type PayOrderRequest struct {
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
	PayType        string `json:"pay_type"`
	Password       string `json:"password"`
	FaceImageURL   string `json:"face_image_url"`
}

type PayOrderResult struct {
	UsedPoints  int     `json:"used_points"`
	UsedBalance float64 `json:"used_balance"`
}

// PayOrder processes a wallet payment for an order with idempotency and concurrency safety.
//
// Flow:
//  1. Idempotency check — if idempotency_key already exists, return cached result
//  2. Insert payment record (init status) with UNIQUE idempotency_key
//  3. Transaction: lock order (FOR UPDATE) → lock wallet (FOR UPDATE) → debit → mark order paid
//  4. Update payment record status
//  5. Best-effort publish order.paid event
func (s *PaymentService) PayOrder(orderID, userID int64, req PayOrderRequest) (*PayOrderResult, error) {
	if req.IdempotencyKey == "" {
		return nil, errors.New("支付参数不完整，请重试")
	}

	// Defer cleanup of the temporary pay face image if one was uploaded
	if req.PayType == "face" && req.FaceImageURL != "" {
		defer s.eventBus.PublishFileCleanup(req.FaceImageURL)
	}

	if err := s.ValidatePayAuth(userID, req); err != nil {
		return nil, err
	}

	// Step 1: Idempotency check
	existing, err := s.paymentRepo.FindByIdempotencyKey(req.IdempotencyKey)
	if err == nil {
		switch existing.Status {
		case consts.PaymentStatusSuccess:
			return &PayOrderResult{}, nil // Already paid
		case consts.PaymentStatusInit:
			// If the init record is older than 2 minutes, it's likely stale/abandoned
			// Mark it as failed and allow retry
			if time.Since(existing.CreatedAt) > 2*time.Minute {
				_, _ = s.paymentRepo.UpdateStatus(nil, existing.ID, consts.PaymentStatusInit, consts.PaymentStatusFailed, "支付超时自动作废")
				// Fall through to allow retry
			} else {
				return nil, errors.New("支付处理中，请稍后重试")
			}
		case consts.PaymentStatusFailed:
			// Allow retry with same key — delete old failed record
			if err := s.db.Delete(&model.PaymentRecord{}, existing.ID).Error; err != nil {
				return nil, err
			}
		}
	}

	// Verify order exists and belongs to user
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	if order.UserID != userID {
		return nil, errors.New("无权操作此订单")
	}
	if order.Status == consts.OrderStatusPaid {
		// Already paid, return success immediately
		return &PayOrderResult{
			UsedPoints:  order.UsedPoints,
			UsedBalance: float64(order.UsedBalance) / 100,
		}, nil
	}
	if order.Status != consts.OrderStatusPendingPayment {
		return nil, errors.New("订单状态不允许支付")
	}

	// Step 2: Insert payment record (init status)
	paymentMethod := "balance"
	if req.PayType == "face" {
		paymentMethod = "face"
	} else if req.PayType == "password" {
		paymentMethod = "password"
	} else if req.PayType == "nopassword" {
		paymentMethod = "nopassword"
	}

	record := &model.PaymentRecord{
		OrderID:        order.ID,
		OrderNo:        order.OrderNo,
		UserID:         userID,
		Amount:         order.TotalAmount,
		PaymentMethod:  paymentMethod,
		Status:         consts.PaymentStatusInit,
		IdempotencyKey: req.IdempotencyKey,
	}
	if err := s.paymentRepo.Create(record); err != nil {
		// UNIQUE constraint violation means concurrent request with same key
		return nil, errors.New("支付请求重复，请稍后重试")
	}

	// Step 3: Core transaction
	var txErr error
	var finalUsedPoints int
	var finalRemainingAmount int64

	err = s.db.Transaction(func(tx *gorm.DB) error {
		orderRepo := s.orderRepo.WithTx(tx)
		walletRepo := s.walletRepo.WithTx(tx)
		storeProductRepo := s.storeProductRepo.WithTx(tx)
		productRepo := s.productRepo.WithTx(tx)
		paymentRepo := s.paymentRepo.WithTx(tx)

		// Lock order row
		lockedOrder, err := orderRepo.FindByIDForUpdate(tx, orderID)
		if err != nil {
			txErr = errors.New("订单不存在")
			return err
		}
		if lockedOrder.Status == consts.OrderStatusPaid {
			// Already paid, return nil (success) and sync status to order memory
			order.UsedPoints = lockedOrder.UsedPoints
			order.UsedBalance = lockedOrder.UsedBalance
			finalUsedPoints = lockedOrder.UsedPoints
			finalRemainingAmount = lockedOrder.UsedBalance
			affected, err := paymentRepo.UpdateStatus(tx, record.ID, consts.PaymentStatusInit, consts.PaymentStatusSuccess, "")
			if err != nil {
				return err
			}
			if affected == 0 {
				current, err := paymentRepo.FindByIdempotencyKey(req.IdempotencyKey)
				if err != nil || current.Status != consts.PaymentStatusSuccess {
					return errors.New("更新支付记录状态失败")
				}
			}
			return nil
		}
		if lockedOrder.Status != consts.OrderStatusPendingPayment {
			txErr = errors.New("订单状态已变更，无法支付")
			return fmt.Errorf("order status changed: %d", lockedOrder.Status)
		}

		// Check if order has expired — if so, cancel and release stock within this tx
		if lockedOrder.ExpireAt != nil && lockedOrder.ExpireAt.Before(time.Now()) {
			affected, err := orderRepo.MarkAsCancelled(tx, orderID, consts.OrderStatusPendingPayment, "订单已过期")
			if err != nil {
				return err
			}
			if affected > 0 {
				for _, item := range lockedOrder.Items {
					if err := storeProductRepo.RestoreStock(tx, lockedOrder.StoreID, item.ProductID, item.Quantity); err != nil {
						return err
					}
					if _, err := productRepo.RestoreStock(tx, item.ProductID, item.Quantity); err != nil {
						return err
					}
				}
			}
			txErr = errors.New("订单已过期，已自动取消")
			return txErr
		}

		// Lock user row to get and update green_points
		var userRecord model.SysUser
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userID).First(&userRecord).Error; err != nil {
			txErr = errors.New("用户不存在")
			return err
		}

		// Calculate points offset (1 point = 10 cents)
		maxDeductiblePoints := int(order.TotalAmount / 10)
		usedPoints := userRecord.GreenPoints
		if usedPoints > maxDeductiblePoints {
			usedPoints = maxDeductiblePoints
		}

		remainingAmountInCents := order.TotalAmount - int64(usedPoints*10)
		finalUsedPoints = usedPoints
		finalRemainingAmount = remainingAmountInCents

		// Lock and debit wallet row if needed
		var walletBalanceBefore int64
		var walletBalanceAfter int64

		if remainingAmountInCents > 0 {
			wallet, err := walletRepo.FindByUserIDForUpdate(tx, userID)
			if err != nil {
				txErr = errors.New("钱包不存在")
				return err
			}
			walletBalanceBefore = wallet.Balance
			walletBalanceAfter = wallet.Balance - remainingAmountInCents

			affected, err := walletRepo.Debit(tx, userID, remainingAmountInCents)
			if err != nil {
				return err
			}
			if affected == 0 {
				txErr = errors.New("余额不足")
				return fmt.Errorf("insufficient balance: need %d, have %d", remainingAmountInCents, wallet.Balance)
			}
		} else {
			wallet, err := walletRepo.FindByUserIDForUpdate(tx, userID)
			if err == nil && wallet != nil {
				walletBalanceBefore = wallet.Balance
				walletBalanceAfter = wallet.Balance
			}
		}

		// Deduct points from user
		if usedPoints > 0 {
			err = tx.Model(&model.SysUser{}).Where("id = ?", userID).Update("green_points", gorm.Expr("green_points - ?", usedPoints)).Error
			if err != nil {
				return err
			}
		}

		// Record wallet transaction
		remark := "钱包支付"
		if usedPoints > 0 {
			if remainingAmountInCents == 0 {
				remark = "积分支付"
			} else {
				remark = "积分+钱包支付"
			}
		}

		if err := walletRepo.CreateTransaction(tx, &model.WalletTransaction{
			UserID:         userID,
			Type:           consts.WalletTxTypeOrderPay,
			Amount:         -remainingAmountInCents,
			BalanceBefore:  walletBalanceBefore,
			BalanceAfter:   walletBalanceAfter,
			RelatedID:      order.ID,
			BizType:        consts.BizTypeOrderPay,
			BizID:          order.OrderNo,
			IdempotencyKey: strPtr("pay:" + order.OrderNo),
			Remark:         remark,
		}); err != nil {
			return err
		}

		// Mark order as paid (conditional: WHERE status=pending_payment)
		affectedRows, err := orderRepo.MarkAsPaid(tx, orderID, remainingAmountInCents, usedPoints)
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			txErr = errors.New("订单状态已变更，无法支付")
			return fmt.Errorf("MarkAsPaid affected 0 rows")
		}

		// Update order in memory for response
		order.UsedPoints = usedPoints
		order.UsedBalance = remainingAmountInCents

		affected, err := paymentRepo.UpdateStatus(tx, record.ID, consts.PaymentStatusInit, consts.PaymentStatusSuccess, "")
		if err != nil {
			return err
		}
		if affected == 0 {
			current, err := paymentRepo.FindByIdempotencyKey(req.IdempotencyKey)
			if err != nil || current.Status != consts.PaymentStatusSuccess {
				return errors.New("更新支付记录状态失败")
			}
		}

		return nil
	})

	if err != nil {
		failReason := "支付失败"
		if txErr != nil {
			failReason = txErr.Error()
		}
		_, _ = s.paymentRepo.UpdateStatus(nil, record.ID, consts.PaymentStatusInit, consts.PaymentStatusFailed, failReason)
		if txErr != nil {
			return nil, txErr
		}
		return nil, errors.New("支付失败，请重试")
	}

	// Step 4: Best-effort event publish
	if s.eventBus != nil {
		order.Status = consts.OrderStatusPaid
		now := time.Now()
		order.PaidAt = &now
		s.eventBus.PublishOrderPaid(order)
	}

	return &PayOrderResult{
		UsedPoints:  finalUsedPoints,
		UsedBalance: float64(finalRemainingAmount) / 100,
	}, nil
}

func (s *PaymentService) ValidatePayAuth(userID int64, req PayOrderRequest) error {
	payType := req.PayType
	if payType == "nopassword" || payType == "none" {
		return nil
	}
	if payType == "" {
		payType = "password"
	}
	switch payType {
	case "password":
		if req.Password == "" {
			return errors.New("请输入支付密码")
		}
		user, err := s.userRepo.FindByID(userID)
		if err != nil {
			return errors.New("用户不存在")
		}
		if user.Status != 1 {
			return errors.New("用户状态异常")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return errors.New("支付密码错误")
		}
	case "face":
		user, err := s.userRepo.FindByID(userID)
		if err != nil {
			return errors.New("用户不存在")
		}
		if user.Status != 1 {
			return errors.New("用户状态异常")
		}
		if !user.FaceRegistered {
			return errors.New("当前账号未录入人脸")
		}
		if req.FaceImageURL == "" {
			return errors.New("请先完成刷脸验证")
		}
	default:
		return errors.New("不支持的支付方式")
	}
	return nil
}

// PaymentStatusResponse is returned by GetPaymentStatus.
type PaymentStatusResponse struct {
	Status int    `json:"status"` // 0=init, 1=success, 2=failed
	Reason string `json:"reason,omitempty"`
}

// GetPaymentStatus returns the payment status for an order.
func (s *PaymentService) GetPaymentStatus(orderID int64) (*PaymentStatusResponse, error) {
	record, err := s.paymentRepo.FindByOrderID(orderID)
	if err != nil {
		return nil, errors.New("未找到支付记录")
	}
	return &PaymentStatusResponse{
		Status: record.Status,
		Reason: record.FailReason,
	}, nil
}

// RabbitMQ delayed message publisher for order timeout.
// Uses TTL + dead-letter-exchange pattern.

const (
	OrderTimeoutQueue      = "order.timeout"
	OrderTimeoutDLX        = "order.timeout.dlx"
	OrderTimeoutDLXQueue   = "order.timeout.dlx.queue"
	OrderTimeoutRoutingKey = "order.timeout"
)

// PublishOrderTimeout publishes a delayed message that will be delivered after order expires.
func (s *PaymentService) PublishOrderTimeout(orderID int64, delay time.Duration) {
	if s.eventBus == nil || s.eventBus.mq == nil {
		return
	}
	body, _ := json.Marshal(map[string]int64{"order_id": orderID})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// For simplicity, publish to the timeout queue directly.
	// In production, use TTL + DLX for true delayed delivery.
	// Here we just log and rely on the polling fallback.
	if err := s.eventBus.mq.PublishEvent(ctx, OrderTimeoutQueue, body); err != nil {
		s.log.Warn("publish order timeout failed", "order_id", orderID, "error", err)
	}
}
