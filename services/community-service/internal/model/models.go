package model

import "time"

type Notice struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:100;not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Publisher string    `gorm:"size:50;not null;default:''" json:"publisher"`
	ViewCount int64     `gorm:"not null;default:0" json:"view_count"`
	Status    int       `gorm:"not null;default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Notice) TableName() string { return "notices" }

type NoticeViewLog struct {
	ID        int64      `gorm:"primaryKey" json:"id"`
	NoticeID  int64      `gorm:"not null;index:idx_notice_user,unique" json:"notice_id"`
	UserID    int64      `gorm:"not null;index:idx_notice_user,unique" json:"user_id"`
	ViewedAt  time.Time  `gorm:"not null" json:"viewed_at"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (NoticeViewLog) TableName() string { return "notice_view_logs" }

type Visitor struct {
	ID           int64      `gorm:"primaryKey" json:"id"`
	UserID       int64      `gorm:"not null;index" json:"user_id"`
	VisitorName  string     `gorm:"size:50;not null" json:"visitor_name"`
	VisitorPhone string     `gorm:"size:20;not null" json:"visitor_phone"`
	VisitPurpose string     `gorm:"size:255;not null" json:"visit_purpose"`
	ReleaseTime  time.Time  `gorm:"not null" json:"release_time"`
	ValidDate    time.Time  `gorm:"not null" json:"valid_date"`
	Status       int        `gorm:"not null;default:0;index" json:"status"` // 0 pending, 1 approved, 2 rejected
	AuditRemark  string     `gorm:"size:255;not null;default:''" json:"audit_remark"`
	AuditAt      *time.Time `json:"audit_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (Visitor) TableName() string { return "visitors" }

type VisitorAdminView struct {
	ID           int64      `json:"id"`
	UserID       int64      `json:"user_id"`
	UserName     string     `json:"user_name"`
	UserMobile   string     `json:"user_mobile"`
	VisitorName  string     `json:"visitor_name"`
	VisitorPhone string     `json:"visitor_phone"`
	VisitPurpose string     `json:"visit_purpose"`
	ReleaseTime  time.Time  `json:"release_time"`
	ValidDate    time.Time  `json:"valid_date"`
	Status       int        `json:"status"`
	AuditRemark  string     `json:"audit_remark"`
	AuditAt      *time.Time `json:"audit_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type ParkingSpace struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	ParkingNo string    `gorm:"size:50;not null;uniqueIndex" json:"parking_no"`
	Status    int       `gorm:"not null;default:0;index" json:"status"` // 0 free, 1 bound
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ParkingSpace) TableName() string { return "parking_spaces" }

type UserParkingBinding struct {
	ID             int64        `gorm:"primaryKey" json:"id"`
	UserID         int64        `gorm:"not null;index" json:"user_id"`
	ParkingSpaceID int64        `gorm:"not null;index" json:"parking_space_id"`
	CarPlate       string       `gorm:"size:20;not null;default:''" json:"car_plate"`
	Status         int          `gorm:"not null;default:1;index" json:"status"` // 1 active, 0 inactive
	BoundAt        time.Time    `gorm:"not null" json:"bound_at"`
	UnboundAt      *time.Time   `json:"unbound_at"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	ParkingSpace   ParkingSpace `gorm:"foreignKey:ParkingSpaceID" json:"parking_space,omitempty"`
}

func (UserParkingBinding) TableName() string { return "user_parking_bindings" }

type ParkingSpaceAdminView struct {
	ID         int64     `json:"id"`
	ParkingNo  string    `json:"parking_no"`
	Status     int       `json:"status"`
	UserID     int64     `json:"user_id"`
	UserName   string    `json:"user_name"`
	UserMobile string    `json:"user_mobile"`
	CarPlate   string    `json:"car_plate"`
	BindingID  int64     `json:"binding_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PropertyFee struct {
	ID        int64      `gorm:"primaryKey" json:"id"`
	UserID    int64      `gorm:"not null;index" json:"user_id"`
	Month     string     `gorm:"size:20;not null;index" json:"month"`
	Amount    int64      `gorm:"not null" json:"amount"`                 // cents
	Status    int        `gorm:"not null;default:0;index" json:"status"` // 0 unpaid, 1 paid
	DueDate   *time.Time `json:"due_date"`
	PaidAt    *time.Time `json:"paid_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (PropertyFee) TableName() string { return "property_fees" }

type PropertyFeePayment struct {
	ID                  int64      `gorm:"primaryKey" json:"id"`
	PropertyFeeID       int64      `gorm:"not null;index" json:"property_fee_id"`
	UserID              int64      `gorm:"not null;index;uniqueIndex:uk_property_payment_user_idempotency" json:"user_id"`
	Amount              int64      `gorm:"not null" json:"amount"`
	WalletTransactionID int64      `gorm:"not null;default:0" json:"wallet_transaction_id"`
	IdempotencyKey      string     `gorm:"size:64;uniqueIndex:uk_property_payment_user_idempotency" json:"idempotency_key"`
	Status              int        `gorm:"not null;default:1;index" json:"status"` // 1 success
	PaidAt              *time.Time `json:"paid_at"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func (PropertyFeePayment) TableName() string { return "property_fee_payments" }

// ── Workorder models ──

const (
	StatusPending    = 0
	StatusProcessing = 1
	StatusCompleted  = 2
	StatusRejected   = 3

	WorkorderTypeRepair    = "repair"
	WorkorderTypeComplaint = "complaint"
)

type WorkOrder struct {
	ID          int64      `gorm:"primaryKey" json:"id"`
	Type        string     `gorm:"size:20;not null;index:idx_workorders_type_status" json:"type"`
	UserID      int64      `gorm:"not null;index" json:"user_id"`
	Category    string     `gorm:"size:50;not null" json:"category"`
	Description string     `gorm:"type:text;not null" json:"description"`
	Status      int        `gorm:"not null;default:0;index:idx_workorders_type_status" json:"status"`
	Result      string     `gorm:"size:500;not null;default:''" json:"result"`
	ProcessorID int64      `gorm:"not null;default:0" json:"processor_id"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (WorkOrder) TableName() string { return "workorders" }

type WorkOrderAdminView struct {
	ID          int64      `json:"id"`
	Type        string     `json:"type"`
	UserID      int64      `json:"user_id"`
	UserName    string     `json:"user_name"`
	UserMobile  string     `json:"user_mobile"`
	Category    string     `json:"category"`
	Description string     `json:"description"`
	Status      int        `json:"status"`
	Result      string     `json:"result"`
	ProcessorID int64      `json:"processor_id"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type WorkorderLog struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	TargetType string    `gorm:"size:20;not null;index:idx_target" json:"target_type"`
	TargetID   int64     `gorm:"not null;index:idx_target" json:"target_id"`
	FromStatus int       `gorm:"not null;default:-1" json:"from_status"`
	ToStatus   int       `gorm:"not null" json:"to_status"`
	OperatorID int64     `gorm:"not null;default:0" json:"operator_id"`
	Action     string    `gorm:"size:50;not null" json:"action"`
	Remark     string    `gorm:"size:500;not null;default:''" json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
}

func (WorkorderLog) TableName() string { return "workorder_logs" }

// ── Statistics read-only models ──

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

// ── Statistics aggregation result structs ──

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

// ── AI Report ──

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

// ── Community Chat ──

type CommunityMessage struct {
	ID        int64                 `gorm:"primaryKey" json:"id"`
	UserID    int64                 `gorm:"not null;index" json:"user_id"`
	Content   string                `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time             `json:"created_at"`
	User      *CommunityMessageUser `gorm:"-" json:"user,omitempty"`
}

type CommunityMessageUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

func (CommunityMessage) TableName() string { return "community_messages" }
