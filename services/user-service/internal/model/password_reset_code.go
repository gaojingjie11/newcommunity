package model

import "time"

type PasswordResetCode struct {
	ID        int64      `gorm:"primaryKey" json:"id"`
	Mobile    string     `gorm:"type:varchar(20);index" json:"mobile"`
	CodeHash  string     `gorm:"type:varchar(255)" json:"-"`
	ExpiresAt time.Time  `gorm:"column:expires_at" json:"expires_at"`
	UsedAt    *time.Time `gorm:"column:used_at" json:"used_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func (PasswordResetCode) TableName() string { return "password_reset_codes" }
