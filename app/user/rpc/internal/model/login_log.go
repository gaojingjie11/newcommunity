package model

import "time"

type UserLoginLog struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	UserID        int64     `gorm:"index" json:"user_id"`
	Mobile        string    `gorm:"type:varchar(20)" json:"mobile"`
	LoginTime     time.Time `gorm:"column:login_time" json:"login_time"`
	IP            string    `gorm:"type:varchar(64)" json:"ip"`
	UserAgent     string    `gorm:"type:varchar(512)" json:"user_agent"`
	Success       bool      `json:"success"`
	FailureReason string    `gorm:"type:varchar(255)" json:"failure_reason"`
}

func (UserLoginLog) TableName() string { return "user_login_logs" }

type AdminLoginLog struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	AdminUserID   int64     `gorm:"column:admin_user_id;index" json:"admin_user_id"`
	Mobile        string    `gorm:"type:varchar(20)" json:"mobile"`
	LoginTime     time.Time `gorm:"column:login_time" json:"login_time"`
	IP            string    `gorm:"type:varchar(64)" json:"ip"`
	UserAgent     string    `gorm:"type:varchar(512)" json:"user_agent"`
	Success       bool      `json:"success"`
	FailureReason string    `gorm:"type:varchar(255)" json:"failure_reason"`
}

func (AdminLoginLog) TableName() string { return "admin_login_logs" }
