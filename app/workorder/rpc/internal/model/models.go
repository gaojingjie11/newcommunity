package model

import "time"

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
	Type        string     `gorm:"size:20;not null;index:idx_workorders_type_status;index:idx_workorders_created_at,priority:1" json:"type"`
	UserID      int64      `gorm:"not null;index" json:"user_id"`
	Category    string     `gorm:"size:50;not null" json:"category"`
	Description string     `gorm:"type:text;not null" json:"description"`
	Status      int        `gorm:"not null;default:0;index:idx_workorders_type_status" json:"status"`
	Result      string     `gorm:"size:500;not null;default:''" json:"result"`
	ProcessorID int64      `gorm:"not null;default:0" json:"processor_id"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt   time.Time  `gorm:"index:idx_workorders_created_at,priority:2" json:"created_at"`
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
