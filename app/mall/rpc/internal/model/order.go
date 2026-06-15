package model

import "time"

type Order struct {
	ID           int64       `gorm:"primaryKey" json:"id"`
	OrderNo      string      `gorm:"column:order_no;type:varchar(64);uniqueIndex" json:"order_no"`
	UserID       int64       `gorm:"index" json:"user_id"`
	StoreID      int64       `gorm:"index" json:"store_id"`
	TotalAmount  int64       `gorm:"type:bigint;not null;default:0" json:"total_amount"`
	UsedPoints   int         `gorm:"column:used_points;not null;default:0" json:"used_points"`
	UsedBalance  int64       `gorm:"column:used_balance;type:bigint;not null;default:0" json:"used_balance"`
	Status       int         `gorm:"index" json:"status"`
	ExpireAt     *time.Time  `gorm:"column:expire_at;index" json:"expire_at"`
	CancelReason string      `gorm:"type:varchar(255)" json:"cancel_reason,omitempty"`
	CancelledAt  *time.Time  `gorm:"column:cancelled_at" json:"cancelled_at,omitempty"`
	PaidAt       *time.Time  `gorm:"column:paid_at" json:"paid_at"`
	Version      int         `gorm:"not null;default:0" json:"version"`
	CreatedAt    time.Time   `gorm:"index" json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Items        []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	Store        Store       `gorm:"foreignKey:StoreID" json:"store"`
}

func (Order) TableName() string { return "oms_order" }

type OrderItem struct {
	ID              int64   `gorm:"primaryKey" json:"id"`
	OrderID         int64   `gorm:"index" json:"order_id"`
	StoreID         int64   `gorm:"index" json:"store_id"`
	ProductID       int64   `json:"product_id"`
	Price           int64   `gorm:"type:bigint;not null;default:0" json:"price"`
	Quantity        int     `json:"quantity"`
	ProductSnapshot string  `gorm:"type:varchar(512)" json:"product_snapshot,omitempty"`
	Product         Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (OrderItem) TableName() string { return "oms_order_item" }
