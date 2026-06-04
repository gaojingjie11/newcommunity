package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"smartcommunity-microservices/common/mq"
)

type EventBus struct {
	mq  *mq.Client
	log *slog.Logger
}

func NewEventBus(mqClient *mq.Client, log *slog.Logger) *EventBus {
	return &EventBus{mq: mqClient, log: log}
}

func (b *EventBus) Publish(ctx context.Context, eventName string, payload interface{}) string {
	if b == nil || b.mq == nil {
		return "skipped"
	}
	body, err := json.Marshal(payload)
	if err != nil {
		if b.log != nil {
			b.log.Warn("marshal workorder event failed", "event", eventName, "error", err)
		}
		return "failed"
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := b.mq.PublishEvent(ctx, eventName, body); err != nil {
		if b.log != nil {
			b.log.Warn("publish workorder event failed", "event", eventName, "error", err)
		}
		return "failed"
	}
	return "published"
}
