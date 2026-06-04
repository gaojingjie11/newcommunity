package model

import "time"

type SysPermission struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"type:varchar(128);uniqueIndex" json:"code"`
	Name      string    `gorm:"type:varchar(128)" json:"name"`
	Resource  string    `gorm:"type:varchar(64)" json:"resource"`
	Method    string    `gorm:"type:varchar(16)" json:"method"`
	Path      string    `gorm:"type:varchar(255)" json:"path"`
	Type      int       `json:"type"`
	Status    int       `gorm:"not null;default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SysPermission) TableName() string { return "sys_permission" }

type SysRolePermission struct {
	ID           int64 `gorm:"primaryKey" json:"id"`
	RoleID       int64 `gorm:"uniqueIndex:uk_role_permission" json:"role_id"`
	PermissionID int64 `gorm:"uniqueIndex:uk_role_permission" json:"permission_id"`
}

func (SysRolePermission) TableName() string { return "sys_role_permission" }
