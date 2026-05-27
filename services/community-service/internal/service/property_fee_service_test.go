package service

import "testing"

func TestPropertyFeeWalletKeyIsBillScoped(t *testing.T) {
	if got := propertyFeeWalletKey(123); got != "community-fee:123" {
		t.Fatalf("propertyFeeWalletKey = %q, want community-fee:123", got)
	}
}
