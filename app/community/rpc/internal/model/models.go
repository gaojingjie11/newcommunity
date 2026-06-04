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
