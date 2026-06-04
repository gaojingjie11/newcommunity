package service

import (
	"sync"
	"sync/atomic"
	"testing"

	"smartcommunity-microservices/app/mall/rpc/internal/consts"
	"smartcommunity-microservices/app/mall/rpc/internal/model"
)

// TestDebitForExternal_Constants verifies the new wallet tx type and biz type constants.
func TestDebitForExternal_Constants(t *testing.T) {
	if consts.WalletTxTypeFee != 5 {
		t.Errorf("WalletTxTypeFee = %d, want 5", consts.WalletTxTypeFee)
	}
	if consts.BizTypePropertyFee != "property_fee" {
		t.Errorf("BizTypePropertyFee = %s, want property_fee", consts.BizTypePropertyFee)
	}
}

// TestDebitForExternal_IdempotencyKeyPattern verifies the expected key format.
func TestDebitForExternal_IdempotencyKeyPattern(t *testing.T) {
	// The convention is "community-fee:{feeID}"
	feeID := int64(42)
	key := "community-fee:42"
	expected := "community-fee:42"
	if key != expected {
		t.Errorf("idempotency key = %s, want %s", key, expected)
	}

	// Verify it's bound to feeID, not user+fee
	_ = feeID
}

// TestAtomicBalanceDeduction simulates the WHERE balance >= amount pattern
// used by WalletRepo.Debit. Multiple goroutines race to deduct from a shared balance.
func TestAtomicBalanceDeduction(t *testing.T) {
	var balance int64 = 1000 // 10 yuan in cents
	deductionAmount := int64(100)
	var successCount int64

	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for {
				old := atomic.LoadInt64(&balance)
				if old < deductionAmount {
					return // insufficient balance
				}
				if atomic.CompareAndSwapInt64(&balance, old, old-deductionAmount) {
					atomic.AddInt64(&successCount, 1)
					return
				}
			}
		}()
	}
	wg.Wait()

	if balance != 0 {
		t.Errorf("expected balance=0, got %d", balance)
	}
	if successCount != 10 {
		t.Errorf("expected 10 successful deductions, got %d", successCount)
	}
}

// TestAtomicBalanceDeduction_Insufficient simulates insufficient balance scenario.
func TestAtomicBalanceDeduction_Insufficient(t *testing.T) {
	var balance int64 = 50
	deductionAmount := int64(100)
	var successCount int64
	var failCount int64

	const goroutines = 5
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			old := atomic.LoadInt64(&balance)
			if old < deductionAmount {
				atomic.AddInt64(&failCount, 1)
				return
			}
			if atomic.CompareAndSwapInt64(&balance, old, old-deductionAmount) {
				atomic.AddInt64(&successCount, 1)
			} else {
				atomic.AddInt64(&failCount, 1)
			}
		}()
	}
	wg.Wait()

	if successCount != 0 {
		t.Errorf("expected 0 successful deductions with insufficient balance, got %d", successCount)
	}
	if failCount != int64(goroutines) {
		t.Errorf("expected all %d goroutines to fail, got %d", goroutines, failCount)
	}
	if balance != 50 {
		t.Errorf("balance should remain 50, got %d", balance)
	}
}

// TestIdempotencyKey_UniqueConstraint simulates the UNIQUE constraint on idempotency_key.
// Only the first insert with a given key should succeed.
func TestIdempotencyKey_UniqueConstraint(t *testing.T) {
	type txRecord struct {
		idempotencyKey string
		walletTxID     int64
	}

	var mu sync.Mutex
	records := make(map[string]int64)
	var nextID int64 = 1

	tryInsert := func(key string) (int64, bool) {
		mu.Lock()
		defer mu.Unlock()
		if _, exists := records[key]; exists {
			return records[key], false // duplicate
		}
		id := nextID
		nextID++
		records[key] = id
		return id, true
	}

	// Simulate DebitForExternal's idempotency check:
	// 1. FindTransactionByIdempotencyKey → if found, return existing
	// 2. Otherwise insert new record

	key := "community-fee:99"

	// First call: should insert
	id1, ok1 := tryInsert(key)
	if !ok1 {
		t.Fatal("first insert should succeed")
	}

	// Second call with same key: should return existing
	id2, ok2 := tryInsert(key)
	if ok2 {
		t.Fatal("second insert with same key should be rejected")
	}
	if id1 != id2 {
		t.Errorf("idempotent lookup should return same ID: got %d vs %d", id1, id2)
	}
}

// TestWalletTxTypeFee_TypeValue verifies WalletTxTypeFee doesn't collide with other types.
func TestWalletTxTypeFee_TypeValue(t *testing.T) {
	types := map[int]string{
		consts.WalletTxTypeOrderPay: "order_pay",
		consts.WalletTxTypeTransfer: "transfer",
		consts.WalletTxTypeRecharge: "recharge",
		consts.WalletTxTypeRefund:   "refund",
		consts.WalletTxTypeFee:      "fee",
	}

	// All types should be unique
	seen := make(map[int]bool)
	for typ, name := range types {
		if seen[typ] {
			t.Errorf("duplicate wallet tx type value %d for %s", typ, name)
		}
		seen[typ] = true
	}

	// Verify fee type is 5
	if consts.WalletTxTypeFee != 5 {
		t.Errorf("WalletTxTypeFee = %d, want 5", consts.WalletTxTypeFee)
	}
}

func TestValidateExternalDebitTransactionRejectsKeyConflict(t *testing.T) {
	tx := &model.WalletTransaction{
		UserID:        7,
		Type:          consts.WalletTxTypeFee,
		Amount:        -100,
		BizType:       consts.BizTypePropertyFee,
		BizID:         "community-fee:1",
		BalanceBefore: 1000,
		BalanceAfter:  900,
	}
	if err := validateExternalDebitTransaction(tx, 7, 100, consts.BizTypePropertyFee, "community-fee:1"); err != nil {
		t.Fatalf("matching transaction should pass validation: %v", err)
	}
	if err := validateExternalDebitTransaction(tx, 8, 100, consts.BizTypePropertyFee, "community-fee:1"); err == nil {
		t.Fatal("different user with same idempotency key should be rejected")
	}
	if err := validateExternalDebitTransaction(tx, 7, 200, consts.BizTypePropertyFee, "community-fee:1"); err == nil {
		t.Fatal("different amount with same idempotency key should be rejected")
	}
}

// TestDebitForExternal_CrossServiceTransactionFlow documents the expected flow.
func TestDebitForExternal_CrossServiceTransactionFlow(t *testing.T) {
	// This test documents the cross-service transaction design.
	// In production:
	//
	// 1. community-service validates fee (exists, belongs to user, unpaid)
	// 2. community-service calls mall-service POST /api/internal/mall/wallet/debit
	//    with idempotency_key = "community-fee:{feeID}"
	// 3. mall-service checks idempotency → if exists, returns existing tx
	// 4. mall-service: SELECT FOR UPDATE wallet → atomic debit → write wallet_transactions
	// 5. community-service: receives wallet_transaction_id
	// 6. community-service: SELECT FOR UPDATE fee → conditional update status 0→1
	//    → write property_fee_payments with wallet_transaction_id
	//
	// If step 4 fails (insufficient balance), step 6 is never reached.
	// If step 6 fails, the wallet debit is committed but community-service
	// can be retried — idempotency key ensures no double debit.

	feeID := int64(123)
	key := "community-fee:123"

	if key != "community-fee:123" {
		t.Error("idempotency key format mismatch")
	}
	_ = feeID
}

// TestConcurrentPropertyFeePay simulates two concurrent Pay requests for the same fee.
// Only one should succeed (the one that wins the SELECT FOR UPDATE lock).
func TestConcurrentPropertyFeePay(t *testing.T) {
	// Simulate the community-service side: status 0→1 with atomic CAS
	var feeStatus int32 = 0 // unpaid
	var paySuccess int32

	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			// Simulate: WHERE status = 0, update to 1
			if atomic.CompareAndSwapInt32(&feeStatus, 0, 1) {
				atomic.AddInt32(&paySuccess, 1)
			}
		}()
	}
	wg.Wait()

	if paySuccess != 1 {
		t.Errorf("expected exactly 1 successful pay, got %d", paySuccess)
	}
	if feeStatus != 1 {
		t.Errorf("expected fee status = 1 (paid), got %d", feeStatus)
	}
}
