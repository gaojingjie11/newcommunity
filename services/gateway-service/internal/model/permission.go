package model

// SysUserRole is a minimal model for gateway permission queries.
type SysUserRole struct {
	ID     int64 `gorm:"primaryKey" json:"id"`
	UserID int64 `gorm:"column:user_id;index" json:"user_id"`
	RoleID int64 `gorm:"column:role_id;index" json:"role_id"`
}

func (SysUserRole) TableName() string { return "sys_user_role" }

// SysRolePermission is a minimal model for gateway permission queries.
type SysRolePermission struct {
	ID           int64 `gorm:"primaryKey" json:"id"`
	RoleID       int64 `gorm:"column:role_id;index" json:"role_id"`
	PermissionID int64 `gorm:"column:permission_id;index" json:"permission_id"`
}

func (SysRolePermission) TableName() string { return "sys_role_permission" }

// SysPermission is a minimal model for gateway permission queries.
type SysPermission struct {
	ID     int64  `gorm:"primaryKey" json:"id"`
	Code   string `gorm:"type:varchar(100);uniqueIndex" json:"code"`
	Status int    `gorm:"not null;default:1" json:"status"`
}

func (SysPermission) TableName() string { return "sys_permission" }
