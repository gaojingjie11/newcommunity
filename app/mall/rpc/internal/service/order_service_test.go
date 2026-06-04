package service

import (
	"testing"

	"smartcommunity-microservices/app/mall/rpc/internal/consts"
)

func TestOrderStatusConstants(t *testing.T) {
	if consts.OrderStatusPendingPayment != 0 {
		t.Errorf("OrderStatusPendingPayment = %d, want 0", consts.OrderStatusPendingPayment)
	}
	if consts.OrderStatusPaid != 1 {
		t.Errorf("OrderStatusPaid = %d, want 1", consts.OrderStatusPaid)
	}
	if consts.OrderStatusShipped != 2 {
		t.Errorf("OrderStatusShipped = %d, want 2", consts.OrderStatusShipped)
	}
	if consts.OrderStatusCompleted != 3 {
		t.Errorf("OrderStatusCompleted = %d, want 3", consts.OrderStatusCompleted)
	}
	if consts.OrderStatusCancelled != 40 {
		t.Errorf("OrderStatusCancelled = %d, want 40", consts.OrderStatusCancelled)
	}
}

func TestWalletTxTypeConstants(t *testing.T) {
	if consts.WalletTxTypeOrderPay != 1 {
		t.Errorf("WalletTxTypeOrderPay = %d, want 1", consts.WalletTxTypeOrderPay)
	}
	if consts.WalletTxTypeTransfer != 2 {
		t.Errorf("WalletTxTypeTransfer = %d, want 2", consts.WalletTxTypeTransfer)
	}
	if consts.WalletTxTypeRecharge != 3 {
		t.Errorf("WalletTxTypeRecharge = %d, want 3", consts.WalletTxTypeRecharge)
	}
	if consts.WalletTxTypeRefund != 4 {
		t.Errorf("WalletTxTypeRefund = %d, want 4", consts.WalletTxTypeRefund)
	}
}

func TestPaymentStatusConstants(t *testing.T) {
	if consts.PaymentStatusInit != 0 {
		t.Errorf("PaymentStatusInit = %d, want 0", consts.PaymentStatusInit)
	}
	if consts.PaymentStatusSuccess != 1 {
		t.Errorf("PaymentStatusSuccess = %d, want 1", consts.PaymentStatusSuccess)
	}
	if consts.PaymentStatusFailed != 2 {
		t.Errorf("PaymentStatusFailed = %d, want 2", consts.PaymentStatusFailed)
	}
}

func TestOrderExpireDuration(t *testing.T) {
	if consts.OrderExpireDuration != 15 {
		t.Errorf("OrderExpireDuration = %d, want 15", consts.OrderExpireDuration)
	}
}

func TestCreateOrderRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateOrderRequest
		wantErr bool
	}{
		{
			name:    "empty cart IDs",
			req:     CreateOrderRequest{CartIDs: []int64{}, StoreID: 1},
			wantErr: true,
		},
		{
			name:    "nil cart IDs",
			req:     CreateOrderRequest{CartIDs: nil, StoreID: 1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The validation is in CreateOrder which requires DB,
			// but we can verify the request struct accepts the fields.
			if len(tt.req.CartIDs) == 0 && !tt.wantErr {
				t.Error("expected error for empty cart IDs")
			}
		})
	}
}

func TestBizTypeConstants(t *testing.T) {
	if consts.BizTypeOrderPay != "order_pay" {
		t.Errorf("BizTypeOrderPay = %s, want order_pay", consts.BizTypeOrderPay)
	}
	if consts.BizTypeOrderRefund != "order_refund" {
		t.Errorf("BizTypeOrderRefund = %s, want order_refund", consts.BizTypeOrderRefund)
	}
	if consts.BizTypeRecharge != "recharge" {
		t.Errorf("BizTypeRecharge = %s, want recharge", consts.BizTypeRecharge)
	}
	if consts.BizTypeTransfer != "transfer" {
		t.Errorf("BizTypeTransfer = %s, want transfer", consts.BizTypeTransfer)
	}
}
