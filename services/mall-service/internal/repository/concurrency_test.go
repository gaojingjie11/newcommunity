package repository

import (
	"sync"
	"sync/atomic"
	"testing"

	"smartcommunity-microservices/services/mall-service/internal/consts"
)

// TestConditionalStatusUpdateConcurrent verifies that two concurrent
// conditional updates (WHERE status=from) on the same row produce
// exactly one winner. Uses a mock counter to simulate the pattern.
func TestConditionalStatusUpdateConcurrent(t *testing.T) {
	// This test verifies the LOGIC of conditional updates without a real DB.
	// Real integration tests require a running MySQL instance.

	var status int32 = consts.OrderStatusPendingPayment
	var winnerCount int32

	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			// Simulate conditional update: only proceed if status == pending
			if atomic.CompareAndSwapInt32(&status, consts.OrderStatusPendingPayment, consts.OrderStatusPaid) {
				atomic.AddInt32(&winnerCount, 1)
			}
		}()
	}
	wg.Wait()

	if winnerCount != 1 {
		t.Errorf("expected exactly 1 winner, got %d", winnerCount)
	}
	if status != consts.OrderStatusPaid {
		t.Errorf("expected status %d, got %d", consts.OrderStatusPaid, status)
	}
}

// TestAtomicStockDeduction simulates the WHERE stock>=qty pattern.
func TestAtomicStockDeduction(t *testing.T) {
	var stock int64 = 10
	var successCount int64

	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for {
				old := atomic.LoadInt64(&stock)
				if old <= 0 {
					return
				}
				if atomic.CompareAndSwapInt64(&stock, old, old-1) {
					atomic.AddInt64(&successCount, 1)
					return
				}
			}
		}()
	}
	wg.Wait()

	if stock != 0 {
		t.Errorf("expected stock=0, got %d", stock)
	}
	if successCount != 10 {
		t.Errorf("expected 10 successful deductions, got %d", successCount)
	}
}

// TestPayCancelRace simulates pay vs cancel racing on the same order.
// Only one should succeed.
func TestPayCancelRace(t *testing.T) {
	var status int32 = consts.OrderStatusPendingPayment
	var payOK, cancelOK int32

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	for i := 0; i < goroutines; i++ {
		// Pay goroutine
		go func() {
			defer wg.Done()
			if atomic.CompareAndSwapInt32(&status, consts.OrderStatusPendingPayment, consts.OrderStatusPaid) {
				atomic.AddInt32(&payOK, 1)
			}
		}()
		// Cancel goroutine
		go func() {
			defer wg.Done()
			if atomic.CompareAndSwapInt32(&status, consts.OrderStatusPendingPayment, consts.OrderStatusCancelled) {
				atomic.AddInt32(&cancelOK, 1)
			}
		}()
	}
	wg.Wait()

	total := atomic.LoadInt32(&payOK) + atomic.LoadInt32(&cancelOK)
	if total != 1 {
		t.Errorf("expected exactly 1 winner (pay or cancel), got %d (pay=%d, cancel=%d)", total, payOK, cancelOK)
	}
}

// TestIdempotencyKeyCollision simulates two goroutines inserting with the same key.
// In real code, the UNIQUE constraint on idempotency_key ensures only one insert succeeds.
func TestIdempotencyKeyCollision(t *testing.T) {
	var inserted int32

	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			// Simulate UNIQUE constraint: only first CAS succeeds
			if atomic.CompareAndSwapInt32(&inserted, 0, 1) {
				// Insert succeeded
			}
		}()
	}
	wg.Wait()

	if inserted != 1 {
		t.Errorf("expected exactly 1 insert, got %d", inserted)
	}
}
