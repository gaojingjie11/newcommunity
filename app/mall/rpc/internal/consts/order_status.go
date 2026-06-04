package consts

// Order status state machine:
//
//	pending_payment(0) --> paid_waiting_ship(1) --> shipped_waiting_receive(2) --> completed(3)
//	pending_payment(0) --> cancelled(40)
//
// No reverse transitions allowed.
const (
	OrderStatusPendingPayment = 0
	OrderStatusPaid           = 1 // paid, waiting for merchant shipment
	OrderStatusShipped        = 2 // shipped, waiting for user receipt
	OrderStatusCompleted      = 3
	OrderStatusCancelled      = 40
)

// Wallet transaction types
const (
	WalletTxTypeOrderPay = 1 // debit for order payment
	WalletTxTypeTransfer = 2 // transfer between users
	WalletTxTypeRecharge = 3 // wallet top-up
	WalletTxTypeRefund   = 4 // refund on cancellation
	WalletTxTypeFee      = 5 // debit for property fee payment
)

// Biz types for wallet_transactions.biz_type
const (
	BizTypeOrderPay    = "order_pay"
	BizTypeOrderRefund = "order_refund"
	BizTypeRecharge    = "recharge"
	BizTypeTransfer    = "transfer"
	BizTypePropertyFee = "property_fee"
)

// Payment record statuses
const (
	PaymentStatusInit    = 0
	PaymentStatusSuccess = 1
	PaymentStatusFailed  = 2
)

// Order timeout
const (
	OrderExpireDuration = 15 // minutes
)
