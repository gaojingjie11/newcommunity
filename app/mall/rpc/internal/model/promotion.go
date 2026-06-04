package model

import "time"

type Promotion struct {
	ID        int64              `gorm:"primaryKey" json:"id"`
	Title     string             `gorm:"type:varchar(128)" json:"title"`
	Type      int                `json:"type"`
	StartDate time.Time          `json:"start_date"`
	EndDate   time.Time          `json:"end_date"`
	Status    int                `json:"status"`
	Products  []PromotionProduct `gorm:"foreignKey:PromotionID" json:"products,omitempty"`
}

func (Promotion) TableName() string { return "pms_promotion" }

type PromotionProduct struct {
	ID          int64 `gorm:"primaryKey" json:"id"`
	PromotionID int64 `gorm:"index" json:"promotion_id"`
	ProductID   int64 `gorm:"index" json:"product_id"`
}

func (PromotionProduct) TableName() string { return "pms_promotion_product" }
