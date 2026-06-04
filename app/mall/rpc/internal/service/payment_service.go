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
			return nil, errors.New("支付处理中，请稍后重试")
		case consts.PaymentStatusFailed:
			// Allow retry with same key — delete old failed record
			// (or could return error; here we allow retry)
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
	if order.Status != consts.OrderStatusPendingPayment {
		return nil, errors.New("订单状态不允许支付")
	}

	// Step 2: Insert payment record (init status)
	record := &model.PaymentRecord{
		OrderID:        order.ID,
		OrderNo:        order.OrderNo,
		UserID:         userID,
		Amount:         order.TotalAmount,
		PaymentMethod:  "wallet",
		Status:         consts.PaymentStatusInit,
		IdempotencyKey: req.IdempotencyKey,
	}
	if err := s.paymentRepo.Create(record); err != nil {
		// UNIQUE constraint violation means concurrent request with same key
		return nil, errors.New("支付请求重复，请稍后重试")
	}

	// Step 3: Core transaction
	var txErr error
	err = s.db.Transaction(func(tx *gorm.DB) error {
		orderRepo := s.orderRepo.WithTx(tx)
		walletRepo := s.walletRepo.WithTx(tx)
		storeProductRepo := s.storeProductRepo.WithTx(tx)
		productRepo := s.productRepo.WithTx(tx)

		// Lock order row
		lockedOrder, err := orderRepo.FindByIDForUpdate(tx, orderID)
		if err != nil {
			txErr = errors.New("订单不存在")
			return err
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

		// Lock wallet row
		wallet, err := walletRepo.FindByUserIDForUpdate(tx, userID)
		if err != nil {
			txErr = errors.New("钱包不存在")
			return err
		}

		// Atomic debit
		affected, err := walletRepo.Debit(tx, userID, order.TotalAmount)
		if err != nil {
			return err
		}
		if affected == 0 {
			txErr = errors.New("余额不足")
			return fmt.Errorf("insufficient balance: need %d, have %d", order.TotalAmount, wallet.Balance)
		}

		// Record wallet transaction
		remark := "钱包支付"
		if order.UsedPoints > 0 {
			if order.UsedBalance == 0 {
				remark = "积分支付"
			} else {
				remark = "积分+钱包支付"
			}
		}
		if err := walletRepo.CreateTransaction(tx, &model.WalletTransaction{
			UserID:         userID,
			Type:           consts.WalletTxTypeOrderPay,
			Amount:         -order.TotalAmount,
			BalanceBefore:  wallet.Balance,
			BalanceAfter:   wallet.Balance - order.TotalAmount,
			RelatedID:      order.ID,
			BizType:        consts.BizTypeOrderPay,
			BizID:          order.OrderNo,
			IdempotencyKey: strPtr("pay:" + order.OrderNo),
			Remark:         remark,
		}); err != nil {
			return err
		}

		// Mark order as paid (conditional: WHERE status=pending_payment)
		affected, err = orderRepo.MarkAsPaid(tx, orderID, order.TotalAmount)
		if err != nil {
			return err
		}
		if affected == 0 {
			txErr = errors.New("订单状态已变更，无法支付")
			return fmt.Errorf("MarkAsPaid affected 0 rows")
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

	// Step 4: Update payment record to success
	_, _ = s.paymentRepo.UpdateStatus(nil, record.ID, consts.PaymentStatusInit, consts.PaymentStatusSuccess, "")

	// Step 5: Best-effort event publish
	if s.eventBus != nil {
		order.Status = consts.OrderStatusPaid
		now := time.Now()
		order.PaidAt = &now
		s.eventBus.PublishOrderPaid(order)
	}

	return &PayOrderResult{
		UsedPoints:  order.UsedPoints,
		UsedBalance: float64(order.TotalAmount) / 100,
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
