package model

type Store struct {
	ID            int64        `gorm:"primaryKey" json:"id"`
	Name          string       `gorm:"type:varchar(128)" json:"name"`
	Address       string       `gorm:"type:varchar(255)" json:"address"`
	Phone         string       `gorm:"type:varchar(32)" json:"phone"`
	AreaID        int64        `gorm:"column:area_id;index" json:"area_id"`
	Region        string       `gorm:"type:varchar(128)" json:"region"`
	BusinessHours string       `gorm:"column:business_hours;type:varchar(64)" json:"business_hours"`
}

func (Store) TableName() string { return "pms_store" }

type StoreProduct struct {
	ID          int64    `gorm:"primaryKey" json:"id"`
	StoreID     int64    `gorm:"index:idx_pms_store_product_store_prod,priority:1" json:"store_id"`
	ProductID   int64    `gorm:"index:idx_pms_store_product_store_prod,priority:2;index" json:"product_id"`
	Stock       int      `json:"stock"`
	LockedStock int      `gorm:"not null;default:0" json:"locked_stock"`
	SoldCount   int      `gorm:"not null;default:0" json:"sold_count"`
	Version     int      `gorm:"not null;default:0" json:"version"`
	Status      int      `gorm:"not null;default:1" json:"status"`
	Product     *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (StoreProduct) TableName() string { return "pms_store_product" }
