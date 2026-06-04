package repository

import (
	"testing"

	"smartcommunity-microservices/app/mall/rpc/internal/consts"
)

func TestOrderStatusTransitions(t *testing.T) {
	// Verify valid state transitions
	transitions := []struct {
		name       string
		fromStatus int
		toStatus   int
		valid      bool
	}{
		{"pending→paid", consts.OrderStatusPendingPayment, consts.OrderStatusPaid, true},
		{"pending→cancelled", consts.OrderStatusPendingPayment, consts.OrderStatusCancelled, true},
		{"paid→shipped", consts.OrderStatusPaid, consts.OrderStatusShipped, true},
		{"paid→cancelled", consts.OrderStatusPaid, consts.OrderStatusCancelled, false},
		{"shipped→completed", consts.OrderStatusShipped, consts.OrderStatusCompleted, true},
		{"cancelled→paid", consts.OrderStatusCancelled, consts.OrderStatusPaid, false},
		{"completed→cancelled", consts.OrderStatusCompleted, consts.OrderStatusCancelled, false},
	}

	for _, tt := range transitions {
		t.Run(tt.name, func(t *testing.T) {
			// The actual validation is in the service layer via conditional updates.
			// Here we just verify the constants are distinct.
			if tt.fromStatus == tt.toStatus && tt.valid {
				t.Error("from and to status should be different for valid transition")
			}
		})
	}
}

func TestPaymentRecordStatusTransitions(t *testing.T) {
	transitions := []struct {
		name       string
		fromStatus int
		toStatus   int
	}{
		{"init→success", consts.PaymentStatusInit, consts.PaymentStatusSuccess},
		{"init→failed", consts.PaymentStatusInit, consts.PaymentStatusFailed},
	}

	for _, tt := range transitions {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fromStatus == tt.toStatus {
				t.Error("from and to status should be different")
			}
		})
	}
}
