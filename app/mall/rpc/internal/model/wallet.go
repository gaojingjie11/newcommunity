package model

import "time"

type Wallet struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `gorm:"uniqueIndex" json:"user_id"`
	Balance   int64     `gorm:"type:bigint;not null;default:0" json:"balance"`
	Version   int       `gorm:"not null;default:0" json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Wallet) TableName() string { return "wallets" }

type WalletTransaction struct {
	ID             int64     `gorm:"primaryKey" json:"id"`
	UserID         int64     `gorm:"index" json:"user_id"`
	Type           int       `json:"type"`
	Amount         int64     `gorm:"type:bigint;not null;default:0" json:"amount"`
	BalanceBefore  int64     `gorm:"type:bigint;not null;default:0" json:"balance_before"`
	BalanceAfter   int64     `gorm:"type:bigint;not null;default:0" json:"balance_after"`
	RelatedID      int64     `json:"related_id"`
	Remark         string    `gorm:"type:varchar(255)" json:"remark"`
	BizType        string    `gorm:"type:varchar(32)" json:"biz_type"`
	BizID          string    `gorm:"type:varchar(64)" json:"biz_id"`
	IdempotencyKey *string   `gorm:"type:varchar(64);uniqueIndex" json:"idempotency_key,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

func (WalletTransaction) TableName() string { return "wallet_transactions" }
