package model

import "time"

type PaymentRecord struct {
	ID             int64      `gorm:"primaryKey" json:"id"`
	OrderID        int64      `gorm:"index;not null" json:"order_id"`
	OrderNo        string     `gorm:"type:varchar(64);not null" json:"order_no"`
	UserID         int64      `gorm:"index;not null" json:"user_id"`
	Amount         int64      `gorm:"type:bigint;not null" json:"amount"`
	PaymentMethod  string     `gorm:"type:varchar(32);not null" json:"payment_method"`
	Status         int        `gorm:"not null;default:0" json:"status"`
	IdempotencyKey string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"idempotency_key"`
	FailReason     string     `gorm:"type:varchar(255)" json:"fail_reason,omitempty"`
	PaidAt         *time.Time `gorm:"column:paid_at" json:"paid_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (PaymentRecord) TableName() string { return "payment_records" }
