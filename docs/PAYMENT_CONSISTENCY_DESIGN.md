# Mall-Service: Payment Consistency Design

## Overview

This document describes the concurrency control and consistency mechanisms implemented in the mall-service order/payment system.

## Key Guarantees

1. **Stock Deduction Safety**: When the last item of a product is purchased concurrently by two users, exactly one succeeds.
2. **Order Expiration**: Orders expire after 30 minutes, automatically cancelling and releasing stock.
3. **Payment Idempotency**: Duplicate payment requests for the same order are idempotent — no double charges.
4. **Cancel/Payment Race**: If cancel and payment race, only one wins based on conditional status update.

---

## Architecture

### State Machine

```
pending_payment(0) ──→ paid(1) ──→ shipped(2) ──→ completed(3)
        │
        └──→ cancelled(40)
```

All transitions use conditional updates: `WHERE id=? AND status=fromStatus`.
`RowsAffected=0` means the state already changed; return appropriate response.

### Components

```
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│ OrderService  │    │PaymentService│    │TimeoutService│
│  CreateOrder  │    │   PayOrder   │    │  Polling     │
│  CancelOrder  │    │   Idempotent │    │  60s interval│
│  ShipOrder    │    │              │    │              │
└──────┬───────┘    └──────┬───────┘    └──────┬───────┘
       │                   │                   │
       ▼                   ▼                   ▼
┌──────────────────────────────────────────────────────┐
│                  MySQL Transactions                   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐           │
│  │  Orders   │  │ Products │  │ Wallets  │           │
│  │  Lock     │  │  Stock   │  │  Balance │           │
│  └──────────┘  └──────────┘  └──────────┘           │
└──────────────────────────────────────────────────────┘
       │
       ▼ (best-effort, after commit)
┌──────────────┐
│   EventBus   │──→ RabbitMQ (order.created, order.paid, order.cancelled)
└──────────────┘
```

---

## Concurrency Control Mechanisms

### 1. Atomic Stock Deduction

```sql
UPDATE products
SET stock = stock - ?, sales = sales + ?, version = version + 1
WHERE id = ? AND stock >= ?
```

- `RowsAffected = 0` → insufficient stock
- Two concurrent requests for the last item: only one gets `RowsAffected = 1`
- No SELECT FOR UPDATE needed for stock; the WHERE clause provides atomicity

### 2. Payment: SELECT FOR UPDATE

Transaction order (fixed to prevent deadlocks):

1. `SELECT ... FROM oms_order WHERE id=? FOR UPDATE` (lock order row)
2. `SELECT ... FROM wallets WHERE user_id=? FOR UPDATE` (lock wallet row)
3. `UPDATE wallets SET balance=balance-? WHERE user_id=? AND balance>=?` (atomic debit)
4. `UPDATE oms_order SET status=1 WHERE id=? AND status=0` (conditional update)

Two concurrent payments on the same order:
- Payment A acquires order lock, Payment B waits
- Payment A commits → status becomes 1
- Payment B acquires lock, sees `status != 0`, returns error

### 3. Idempotency Key

```
Client generates idempotency_key (UUID or business key)
        │
        ▼
┌─────────────────────────┐
│ payment_records table   │
│ idempotency_key UNIQUE  │
└─────────────────────────┘
```

Flow:
1. Check existing record by idempotency_key
2. If `status = success` → return success (idempotent)
3. If `status = init` → return "processing"
4. If not found → insert init record → execute transaction → update status

UNIQUE constraint prevents duplicate inserts from concurrent requests.

### 4. Order Expiration

**Dual Strategy**:
- **Primary**: RabbitMQ delayed message (TTL + DLX) — not yet implemented
- **Fallback**: Polling every 60 seconds: `SELECT * FROM oms_order WHERE status=0 AND expire_at < NOW()`

Both converge to the same `processExpiredOrder`:
1. Conditional cancel: `WHERE status = 0`
2. Restore stock for each order item
3. Best-effort publish `order.cancelled` event

---

## Transaction Boundaries

### CreateOrder

```
BEGIN
  DeductStock (product) × N items
  CreateOrder (order + items)
  DeleteByIDs (cart)
COMMIT
→ Best-effort: PublishOrderCreated
```

### PayOrder

```
BEGIN
  FindByIDForUpdate (order lock)
  FindByUserIDForUpdate (wallet lock)
  Debit (wallet, atomic)
  CreateTransaction (wallet_tx)
  MarkAsPaid (conditional: status=0→1)
COMMIT
→ Update payment_record status
→ Best-effort: PublishOrderPaid
```

### CancelOrder (pending)

```
BEGIN
  MarkAsCancelled (conditional: status=0→40)
  RestoreStock × N items
COMMIT
→ Best-effort: PublishOrderCancelled
```

### CancelOrder (paid)

```
BEGIN
  MarkAsCancelled (conditional: status=1→40)
  RestoreStock × N items
  Credit (wallet refund)
  CreateTransaction (wallet_tx, refund type)
COMMIT
→ Best-effort: PublishOrderCancelled
```

---

## Amount Representation

All monetary values are stored as `int64` in units of **cents** (分).

| Field | Type | Unit |
|-------|------|------|
| Product.Price | int64 | cents |
| Product.OriginalPrice | int64 | cents |
| Order.TotalAmount | int64 | cents |
| Order.UsedBalance | int64 | cents |
| Wallet.Balance | int64 | cents |
| WalletTransaction.Amount | int64 | cents |
| PaymentRecord.Amount | int64 | cents |

Frontend is responsible for ÷100 display conversion.

---

## Data Model Changes

### New Fields

| Table | Column | Type | Description |
|-------|--------|------|-------------|
| oms_order | expire_at | DATETIME | Order expiration time |
| oms_order | cancel_reason | VARCHAR(255) | Cancel reason |
| oms_order | cancelled_at | DATETIME | Cancel timestamp |
| oms_order | paid_at | DATETIME | Payment timestamp |
| oms_order | version | INT | Optimistic lock version |
| oms_order_item | product_snapshot | VARCHAR(512) | Product info at order time |
| pms_product | version | INT | Optimistic lock version |
| wallets | version | INT | Optimistic lock version |
| wallet_transactions | balance_before | BIGINT | Balance before transaction |
| wallet_transactions | balance_after | BIGINT | Balance after transaction |
| wallet_transactions | biz_type | VARCHAR(32) | Business type |
| wallet_transactions | biz_id | VARCHAR(64) | Business ID |
| wallet_transactions | idempotency_key | VARCHAR(64) | Idempotency key (UNIQUE) |
| store_products | locked_stock | INT | Locked stock count |
| store_products | sold_count | INT | Total sold count |
| store_products | version | INT | Optimistic lock version |

### New Table: payment_records

```sql
CREATE TABLE payment_records (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_id        BIGINT NOT NULL,
    order_no        VARCHAR(64) NOT NULL,
    user_id         BIGINT NOT NULL,
    amount          BIGINT NOT NULL DEFAULT 0,
    payment_method  VARCHAR(32) NOT NULL DEFAULT 'wallet',
    status          INT NOT NULL DEFAULT 0,
    idempotency_key VARCHAR(64) NOT NULL,
    fail_reason     VARCHAR(255) DEFAULT '',
    paid_at         DATETIME DEFAULT NULL,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_idempotency_key (idempotency_key),
    INDEX idx_order_id (order_id),
    INDEX idx_user_id (user_id)
);
```

---

## API Contract Changes

### POST /api/mall/orders/:id/pay

Request body now requires `idempotency_key`:

```json
{
  "idempotency_key": "uuid-or-business-key"
}
```

### POST /api/mall/orders

Response now includes `expire_at`:

```json
{
  "order_id": 123,
  "order_no": "17162688000001",
  "total_amount": 5000,
  "expire_at": "2024-05-21T12:30:00Z"
}
```

### POST /api/mall/orders/:id/cancel (NEW)

```json
{
  "reason": "用户取消"  // optional, defaults to "用户取消"
}
```

### GET /api/mall/orders/:id/payment-status (NEW)

```json
{
  "status": 1,      // 0=init, 1=success, 2=failed
  "reason": ""
}
```

---

## Event Publishing

Events are published **best-effort** after transaction commit. Failures are logged but never block the main flow.

| Event | Queue | Trigger |
|-------|-------|---------|
| order.created | `order.created` | After CreateOrder commit |
| order.paid | `order.paid` | After PayOrder commit |
| order.cancelled | `order.cancelled` | After CancelOrder commit |

---

## Graceful Degradation

- **RabbitMQ unavailable**: Service starts without event publishing. Orders/payments work normally.
- **No idempotency_key**: Request rejected with 400 error.
- **Timeout polling fails**: Next poll cycle (60s) will retry.

---

## Cross-Service: Property Fee Wallet Deduction (COMM-007)

### Overview

community-service's property fee payment now calls mall-service's internal wallet API for real deduction, instead of just marking the fee as paid locally.

### Architecture

```
┌──────────────────┐         HTTP (X-Internal-Token)        ┌──────────────────┐
│ community-service │ ──────────────────────────────────────→ │   mall-service   │
│                   │  POST /api/internal/mall/wallet/debit  │                  │
│  PropertyFeeSvc   │         idempotency_key                │  InternalHandler │
│  Pay()            │  ← wallet_transaction_id ──────────────│  DebitWallet()   │
└──────┬────────────┘                                        └──────┬───────────┘
       │                                                            │
       ▼                                                            ▼
┌──────────────────────┐                               ┌──────────────────────┐
│  smart_community DB  │                               │  smart_community DB  │
│  property_fees       │                               │  wallets             │
│  property_fee_payments│                              │  wallet_transactions │
└──────────────────────┘                               └──────────────────────┘
```

### Transaction Strategy

**"Bill-scoped idempotent debit, then conditional-update"** — two independent transactions across one database:

```
1. community-service: Validate fee (exists, belongs to user, unpaid)    [no lock]
2. community-service: Generate wallet idempotency key community-fee:{feeID}
3. community-service: HTTP → mall-service POST /wallet/debit            [cross-service]
   with X-Internal-Token
4. mall-service: Idempotency check (wallet_transactions.idempotency_key)
   and validate existing tx matches user_id/amount/biz_type/biz_id
5. mall-service: BEGIN → GetOrCreate wallet → SELECT FOR UPDATE
                   → UPDATE wallets SET balance=balance-? WHERE balance>=?
                   → INSERT wallet_transactions → COMMIT
6. community-service: Receives wallet_transaction_id
7. community-service: BEGIN → SELECT FOR UPDATE fee
                   → UPDATE property_fees SET status=1 WHERE status=0
                   → INSERT property_fee_payments → COMMIT
```

### Idempotency Key

The idempotency key for property fee payment is `community-fee:{feeID}`:
- Bound to the fee ID, not user+fee — ensures exactly one wallet transaction per fee
- Generated by community-service and not trusted from the client request body/header
- mall-service checks `wallet_transactions.idempotency_key` UNIQUE + app-level lookup
- mall-service rejects a reused key when user, amount, biz_type or biz_id differs
- community-service checks `property_fee_payments.(user_id, idempotency_key)` composite unique index

### Concurrency Guarantees

| Scenario | Mechanism |
|----------|-----------|
| Two concurrent Pay on same fee | Both requests use the same `community-fee:{feeID}` wallet key, so mall-service returns the same wallet tx instead of double debiting; community-service still uses SELECT FOR UPDATE + `WHERE status=0` so only one local payment succeeds |
| Duplicate idempotency_key to mall-service | App-level lookup returns existing wallet tx (no double debit) |
| Idempotency key reused for another user/amount/business | mall-service rejects as key conflict |
| Insufficient balance | mall-service: `WHERE balance>=amount` → RowsAffected=0 → error → fee stays unpaid |
| mall-service debit succeeds, community-service update fails | Retry returns same wallet tx via idempotency key; community-service can re-attempt local update |

### Files Changed

| Service | File | Change |
|---------|------|--------|
| mall-service | `internal/consts/order_status.go` | Added `WalletTxTypeFee=5`, `BizTypePropertyFee` |
| mall-service | `internal/repository/wallet_repo.go` | Added `FindTransactionByIdempotencyKey` |
| mall-service | `internal/service/wallet_service.go` | Added `DebitForExternal` method |
| mall-service | `internal/handler/internal_handler.go` | Added `DebitWallet` handler |
| mall-service | `internal/router/router.go` | Added `POST /wallet/debit` to internal routes |
| mall-service | `configs/config.yaml(.example)` | Added `internal_token` so internal routes are registered and protected |
| community-service | `internal/client/mall_client.go` | NEW — HTTP client for mall-service internal API |
| community-service | `internal/service/property_fee_service.go` | `Pay` now uses service-generated bill-scoped wallet key before local update |
| community-service | `internal/repository/property_fee_repo.go` | `Pay` accepts `walletTxID`, conditional update |
| community-service | `cmd/server/main.go` | MallClient DI wiring |
| community-service | `configs/config.yaml` | Added `mall_internal` section |

---

## Files Changed

| Phase | File | Type |
|-------|------|------|
| 1 | internal/consts/order_status.go | NEW |
| 1 | internal/model/order.go | MODIFIED |
| 1 | internal/model/product.go | MODIFIED |
| 1 | internal/model/wallet.go | MODIFIED |
| 1 | internal/model/payment_record.go | NEW |
| 1 | internal/model/store.go | MODIFIED |
| 2 | internal/repository/order_repo.go | MODIFIED |
| 2 | internal/repository/product_repo.go | MODIFIED |
| 2 | internal/repository/cart_repo.go | MODIFIED |
| 2 | internal/repository/wallet_repo.go | MODIFIED |
| 2 | internal/repository/payment_repo.go | NEW |
| 2 | internal/repository/store_product_repo.go | MODIFIED |
| 3 | internal/service/order_service.go | REWRITTEN |
| 3 | internal/service/payment_service.go | NEW |
| 3 | internal/service/order_timeout_service.go | NEW |
| 3 | internal/service/order_events.go | NEW |
| 3 | internal/service/wallet_service.go | REWRITTEN |
| 4 | internal/handler/order_handler.go | REWRITTEN |
| 4 | internal/handler/payment_handler.go | NEW |
| 4 | internal/handler/wallet_handler.go | MODIFIED |
| 4 | internal/handler/admin_order_handler.go | MODIFIED |
| 5 | internal/router/router.go | MODIFIED |
| 5 | cmd/server/main.go | MODIFIED |
| 6 | migrations/001_order_payment_overhaul.sql | NEW |
