package model

import "time"

type SysUser struct {
	ID             int64     `gorm:"primaryKey" json:"id"`
	Username       string    `gorm:"type:varchar(64)" json:"username"`
	Password       string    `gorm:"type:varchar(255)" json:"-"`
	RealName       string    `gorm:"column:real_name;type:varchar(64)" json:"real_name"`
	Mobile         string    `gorm:"type:varchar(20);uniqueIndex" json:"mobile"`
	Age            int       `json:"age"`
	Gender         int       `json:"gender"`
	Email          string    `gorm:"type:varchar(128)" json:"email"`
	Avatar         string    `gorm:"type:varchar(255)" json:"avatar"`
	GreenPoints    int       `gorm:"column:green_points;not null;default:0" json:"green_points"`
	Balance        float64   `gorm:"type:decimal(10,2);not null;default:0.00" json:"balance"`
	FaceRegistered bool      `gorm:"column:face_registered;not null;default:false" json:"face_registered"`
	FaceImageURL   string    `gorm:"column:face_image_url;type:varchar(512)" json:"face_image_url"`
	Role           string    `gorm:"type:varchar(32)" json:"role"`
	Status         int       `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (SysUser) TableName() string {
	return "sys_user"
}

type SysRole struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(64)" json:"name"`
	Code      string    `gorm:"type:varchar(64);uniqueIndex" json:"code"`
	Remark    string    `gorm:"type:varchar(255)" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
}

func (SysRole) TableName() string {
	return "sys_role"
}

type SysMenu struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	ParentID  int64     `json:"parent_id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Component string    `json:"component"`
	Sort      int       `json:"sort"`
	Type      int       `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

func (SysMenu) TableName() string {
	return "sys_menu"
}

type SysUserRole struct {
	ID     int64 `gorm:"primaryKey" json:"id"`
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}

func (SysUserRole) TableName() string {
	return "sys_user_role"
}

type SysRoleMenu struct {
	ID     int64 `gorm:"primaryKey" json:"id"`
	RoleID int64 `json:"role_id"`
	MenuID int64 `json:"menu_id"`
}

func (SysRoleMenu) TableName() string {
	return "sys_role_menu"
}
