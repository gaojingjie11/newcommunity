package model

import "time"

type Favorite struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `gorm:"index" json:"user_id"`
	ProductID int64     `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"`
}

func (Favorite) TableName() string { return "pms_favorite" }
