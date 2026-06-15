package model

import "time"

type Cart struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `gorm:"index" json:"user_id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"`
}

func (Cart) TableName() string { return "oms_cart" }

type CartItemParam struct {
	CartID   int64 `json:"cart_id"`
	Quantity int   `json:"quantity"`
}
