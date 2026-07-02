package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/common/mq"
)

type OrderPaidEvent struct {
	Event   string `json:"event"`
	OrderID int64  `json:"order_id"`
	OrderNo string `json:"order_no"`
	UserID  int64  `json:"user_id"`
	Amount  int64  `json:"amount"`
	PaidAt  string `json:"paid_at"`
}

type WalletRechargedEvent struct {
	Event          string `json:"event"`
	UserID         int64  `json:"user_id"`
	Amount         int64  `json:"amount"`
	IdempotencyKey string `json:"idempotency_key"`
	RechargedAt    string `json:"recharged_at"`
}

type OrderCancelledEvent struct {
	Event       string `json:"event"`
	OrderID     int64  `json:"order_id"`
	OrderNo     string `json:"order_no"`
	UserID      int64  `json:"user_id"`
	Reason      string `json:"reason"`
	CancelledAt string `json:"cancelled_at"`
}

const (
	QueueOrderPaid       = "order.paid"
	QueueOrderCancelled  = "order.cancelled"
	QueueOrderDelay      = "order.delay.v2"
	QueueOrderTimeout    = "order.timeout.trigger"
	QueueWalletRecharged = "wallet.recharged"
)

type EventBus struct {
	mq  *mq.Client
	log *slog.Logger
}

func NewEventBus(mqClient *mq.Client, log *slog.Logger) *EventBus {
	return &EventBus{mq: mqClient, log: log}
}

func (b *EventBus) PublishOrderDelayCancel(order *model.Order) {
	if b == nil || b.mq == nil {
		return
	}
	event := OrderCancelledEvent{
		Event:   QueueOrderTimeout,
		OrderID: order.ID,
		OrderNo: order.OrderNo,
		UserID:  order.UserID,
		Reason:  "订单超时自动取消",
	}
	body, err := json.Marshal(event)
	if err != nil {
		b.log.Warn("marshal event failed", "queue", QueueOrderDelay, "error", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	delayMs := order.ExpireAt.Sub(time.Now()).Milliseconds()
	if delayMs < 0 {
		delayMs = 0
	}

	b.log.Info("publishing order delay cancel event", "order_id", order.ID, "delay_ms", delayMs)
	if err := b.mq.PublishDelayEvent(ctx, QueueOrderDelay, QueueOrderTimeout, delayMs, body); err != nil {
		b.log.Warn("publish delay event failed", "queue", QueueOrderDelay, "error", err)
	}
}

func (b *EventBus) PublishOrderPaid(order *model.Order) {
	if b == nil || b.mq == nil {
		return
	}
	event := OrderPaidEvent{
		Event:   QueueOrderPaid,
		OrderID: order.ID,
		OrderNo: order.OrderNo,
		UserID:  order.UserID,
		Amount:  order.TotalAmount,
	}
	if order.PaidAt != nil {
		event.PaidAt = order.PaidAt.Format(time.RFC3339)
	}
	b.publish(QueueOrderPaid, event)
}

func (b *EventBus) PublishOrderCancelled(order *model.Order, reason string) {
	if b == nil || b.mq == nil {
		return
	}
	event := OrderCancelledEvent{
		Event:   QueueOrderCancelled,
		OrderID: order.ID,
		OrderNo: order.OrderNo,
		UserID:  order.UserID,
		Reason:  reason,
	}
	if order.CancelledAt != nil {
		event.CancelledAt = order.CancelledAt.Format(time.RFC3339)
	}
	b.publish(QueueOrderCancelled, event)
}

func (b *EventBus) PublishWalletRecharged(userID int64, amount int64, idempotencyKey string) {
	if b == nil || b.mq == nil {
		return
	}
	event := WalletRechargedEvent{
		Event:          QueueWalletRecharged,
		UserID:         userID,
		Amount:         amount,
		IdempotencyKey: idempotencyKey,
		RechargedAt:    time.Now().Format(time.RFC3339),
	}
	b.publish(QueueWalletRecharged, event)
}

func (b *EventBus) publish(queue string, payload interface{}) {
	body, err := json.Marshal(payload)
	if err != nil {
		b.log.Warn("marshal event failed", "queue", queue, "error", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := b.mq.PublishEvent(ctx, queue, body); err != nil {
		b.log.Warn("publish event failed", "queue", queue, "error", err)
	}
}

func (b *EventBus) PublishFileCleanup(url string) {
	if b == nil || b.mq == nil || url == "" {
		return
	}
	event := map[string]string{"url": url}
	b.publish("file.cleanup", event)
}
