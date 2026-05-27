package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"smartcommunity-microservices/pkg/rabbitmq"
	"smartcommunity-microservices/services/mall-service/internal/model"
)

type OrderCreatedEvent struct {
	Event     string `json:"event"`
	OrderID   int64  `json:"order_id"`
	OrderNo   string `json:"order_no"`
	UserID    int64  `json:"user_id"`
	Amount    int64  `json:"amount"`
	CreatedAt string `json:"created_at"`
}

type OrderPaidEvent struct {
	Event   string `json:"event"`
	OrderID int64  `json:"order_id"`
	OrderNo string `json:"order_no"`
	UserID  int64  `json:"user_id"`
	Amount  int64  `json:"amount"`
	PaidAt  string `json:"paid_at"`
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
	QueueOrderCreated   = "order.created"
	QueueOrderPaid      = "order.paid"
	QueueOrderCancelled = "order.cancelled"
)

type EventBus struct {
	mq  *rabbitmq.Client
	log *slog.Logger
}

func NewEventBus(mq *rabbitmq.Client, log *slog.Logger) *EventBus {
	return &EventBus{mq: mq, log: log}
}

func (b *EventBus) PublishOrderCreated(order *model.Order) {
	if b == nil || b.mq == nil {
		return
	}
	event := OrderCreatedEvent{
		Event:     QueueOrderCreated,
		OrderID:   order.ID,
		OrderNo:   order.OrderNo,
		UserID:    order.UserID,
		Amount:    order.TotalAmount,
		CreatedAt: order.CreatedAt.Format(time.RFC3339),
	}
	b.publish(QueueOrderCreated, event)
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
