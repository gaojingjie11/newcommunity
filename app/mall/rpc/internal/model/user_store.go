package model

import "time"

type UserStore struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `gorm:"column:user_id;uniqueIndex:idx_user_store" json:"user_id"`
	StoreID   int64     `gorm:"column:store_id;uniqueIndex:idx_user_store;index:idx_store_id" json:"store_id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (UserStore) TableName() string { return "pms_user_store" }
