package model

import "time"

type Product struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	CategoryName  string    `gorm:"column:category_name;type:varchar(64)" json:"category_name"`
	Name          string    `gorm:"type:varchar(128)" json:"name"`
	Description   string    `gorm:"type:text" json:"description"`
	Price         int64     `gorm:"type:bigint;not null;default:0" json:"price"`
	OriginalPrice int64     `gorm:"type:bigint;not null;default:0" json:"original_price"`
	Stock         int       `json:"stock"`
	ImageURL      string    `gorm:"column:image_url;type:varchar(255)" json:"image_url"`
	IsPromotion   int       `json:"is_promotion"`
	Sales         int       `json:"sales"`
	Status        int       `gorm:"index" json:"status"`
	Version       int       `gorm:"not null;default:0" json:"version"`
	ViewCount     int64     `gorm:"not null;default:0" json:"view_count"`
	CreatedAt     time.Time `json:"created_at"`
	CategoryID    int64     `gorm:"index" json:"category_id"`
}

func (Product) TableName() string { return "pms_product" }

type ProductCategory struct {
	ID   int64  `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(64)" json:"name"`
	Icon string `gorm:"type:varchar(255)" json:"icon"`
	Sort int    `json:"sort"`
}

func (ProductCategory) TableName() string { return "pms_product_category" }
