package model

import "time"

type ProductViewLog struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	ProductID int64     `gorm:"not null" json:"product_id"`
	UserID    int64     `gorm:"not null;default:0" json:"user_id"`
	IP        string    `gorm:"type:varchar(64);not null;default:''" json:"ip"`
	UserAgent string    `gorm:"type:varchar(512);not null;default:''" json:"user_agent"`
	ViewedAt  time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"viewed_at"`
}

func (ProductViewLog) TableName() string { return "product_view_logs" }
