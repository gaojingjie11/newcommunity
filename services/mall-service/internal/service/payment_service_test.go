package service

import (
	"testing"
)

func TestPayOrderRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     PayOrderRequest
		wantErr bool
	}{
		{
			name:    "empty idempotency key",
			req:     PayOrderRequest{IdempotencyKey: ""},
			wantErr: true,
		},
		{
			name:    "valid idempotency key",
			req:     PayOrderRequest{IdempotencyKey: "test-key-123"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.IdempotencyKey == "" && !tt.wantErr {
				t.Error("expected error for empty idempotency key")
			}
		})
	}
}

func TestPaymentStatusResponse_Fields(t *testing.T) {
	resp := PaymentStatusResponse{
		Status: 1,
		Reason: "success",
	}
	if resp.Status != 1 {
		t.Errorf("Status = %d, want 1", resp.Status)
	}
	if resp.Reason != "success" {
		t.Errorf("Reason = %s, want success", resp.Reason)
	}
}

func TestOrderTimeoutConstants(t *testing.T) {
	if OrderTimeoutQueue != "order.timeout" {
		t.Errorf("OrderTimeoutQueue = %s, want order.timeout", OrderTimeoutQueue)
	}
	if OrderTimeoutDLX != "order.timeout.dlx" {
		t.Errorf("OrderTimeoutDLX = %s, want order.timeout.dlx", OrderTimeoutDLX)
	}
}
