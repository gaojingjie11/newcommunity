package model

import "time"

// Statistics read-only models
type Product struct {
	ID    int64  `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Sales int    `json:"sales"`
}

func (Product) TableName() string { return "pms_product" }

type Order struct {
	ID          int64      `gorm:"primaryKey" json:"id"`
	OrderNo     string     `json:"order_no"`
	UserID      int64      `json:"user_id"`
	TotalAmount int64      `json:"total_amount"`
	Status      int        `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	PaidAt      *time.Time `json:"paid_at"`
}

func (Order) TableName() string { return "oms_order" }

type OrderItem struct {
	ID        int64 `gorm:"primaryKey" json:"id"`
	OrderID   int64 `json:"order_id"`
	ProductID int64 `json:"product_id"`
	Price     int64 `json:"price"`
	Quantity  int   `json:"quantity"`
}

func (OrderItem) TableName() string { return "oms_order_item" }

type SysUserStats struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (SysUserStats) TableName() string { return "sys_user" }

type ProductViewLog struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	ProductID int64     `json:"product_id"`
	UserID    int64     `json:"user_id"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	ViewedAt  time.Time `json:"viewed_at"`
}

func (ProductViewLog) TableName() string { return "product_view_logs" }

// Statistics aggregation result structs
type ProductSalesRank struct {
	ProductID   int64  `json:"product_id"`
	ProductName string `json:"product_name"`
	TotalSales  int64  `json:"total_sales"`
	TotalAmount int64  `json:"total_amount"`
}

type ProductViewRank struct {
	ProductID   int64  `json:"product_id"`
	ProductName string `json:"product_name"`
	ViewCount   int64  `json:"view_count"`
	UniqueUsers int64  `json:"unique_users"`
}

type OrderSummary struct {
	Status      int   `json:"status"`
	Count       int64 `json:"count"`
	TotalAmount int64 `json:"total_amount"`
}

type OrderTrend struct {
	Date   string `json:"date"`
	Count  int64  `json:"count"`
	Amount int64  `json:"amount"`
}

type WorkorderSummary struct {
	Type   string `json:"type"`
	Status int    `json:"status"`
	Count  int64  `json:"count"`
}

type CommunityOverview struct {
	UserCount      int64 `json:"user_count"`
	OrderCount     int64 `json:"order_count"`
	PaidAmount     int64 `json:"paid_amount"`
	RepairCount    int64 `json:"repair_count"`
	ComplaintCount int64 `json:"complaint_count"`
	FeeCount       int64 `json:"fee_count"`
	FeePaidCount   int64 `json:"fee_paid_count"`
}

// AI Report
type AIReport struct {
	ID                 int64     `gorm:"primaryKey" json:"id"`
	RepairNewCount     int64     `gorm:"column:repair_new_count;not null;default:0" json:"repair_new_count"`
	RepairPendingCount int64     `gorm:"column:repair_pending_count;not null;default:0" json:"repair_pending_count"`
	VisitorNewCount    int64     `gorm:"column:visitor_new_count;not null;default:0" json:"visitor_new_count"`
	PropertyPaidCount  int64     `gorm:"column:property_paid_count;not null;default:0" json:"property_paid_count"`
	PropertyPaidAmount float64   `gorm:"column:property_paid_amount;type:decimal(10,2);not null;default:0.00" json:"property_paid_amount"`
	ReportSummary      string    `gorm:"column:report_summary;type:varchar(255)" json:"report_summary"`
	Report             string    `gorm:"column:report_markdown;type:longtext" json:"report"`
	GeneratedBy        int64     `gorm:"column:generated_by;not null;default:0" json:"generated_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (AIReport) TableName() string { return "cms_ai_report" }
