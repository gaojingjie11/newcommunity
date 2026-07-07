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
	ID                     int64     `gorm:"primaryKey" json:"id"`
	Username               string    `json:"username"`
	RealName               string    `json:"real_name"`
	GreenPoints            int64     `json:"green_points"`
	GreenPointsTotalEarned int64     `json:"green_points_total_earned"`
	Status                 int       `json:"status"`
	CreatedAt              time.Time `json:"created_at"`
}

func (SysUserStats) TableName() string { return "sys_user" }

type PropertyFee struct {
	ID        int64      `gorm:"primaryKey" json:"id"`
	UserID    int64      `json:"user_id"`
	Month     string     `json:"month"`
	Amount    int64      `json:"amount"` // cents
	Status    int        `json:"status"` // 0 unpaid, 1 paid
	PaidAt    *time.Time `json:"paid_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func (PropertyFee) TableName() string { return "property_fees" }

type ParkingSpace struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	ParkingNo string    `json:"parking_no"`
	Status    int       `json:"status"` // 0 free, 1 bound
	CreatedAt time.Time `json:"created_at"`
}

func (ParkingSpace) TableName() string { return "parking_spaces" }

type UserParkingBinding struct {
	ID             int64     `gorm:"primaryKey" json:"id"`
	UserID         int64     `json:"user_id"`
	ParkingSpaceID int64     `json:"parking_space_id"`
	CarPlate       string    `json:"car_plate"`
	Status         int       `json:"status"` // 1 active, 0 inactive
	CreatedAt      time.Time `json:"created_at"`
}

func (UserParkingBinding) TableName() string { return "user_parking_bindings" }

type WorkOrder struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Type      string    `json:"type"`
	Category  string    `json:"category"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (WorkOrder) TableName() string { return "workorders" }

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

type RepairStat struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type EcoLeaderboard struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	RealName    string `json:"real_name"`
	GreenPoints int64  `json:"green_points"`
}

type CommunityOverview struct {
	UserCount      int64 `json:"user_count"`
	OrderCount     int64 `json:"order_count"`
	PaidAmount     int64 `json:"paid_amount"`
	RepairCount    int64 `json:"repair_count"`
	ComplaintCount int64 `json:"complaint_count"`
	FeeCount       int64 `json:"fee_count"`
	FeePaidCount   int64 `json:"fee_paid_count"`

	TotalUsers    int64        `json:"total_users"`
	TodayOrders   int64        `json:"today_orders"`
	ParkingRate   string       `json:"parking_rate"`
	MonthIncome   float64      `json:"month_income"`
	RepairStats   []RepairStat `json:"repair_stats"`
	IncomeDates   []string     `json:"income_dates"`
	IncomeTrend   []float64    `json:"income_trend"`
	CostStructure []float64    `json:"cost_structure"`
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
	Report             string    `gorm:"column:report_markdown;type:text" json:"report"`
	GeneratedBy        int64     `gorm:"column:generated_by;not null;default:0" json:"generated_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (AIReport) TableName() string { return "cms_ai_report" }
