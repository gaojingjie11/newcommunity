package model

// SysUser is a lean read-only reference for cross-service queries.
// This table is owned by user-service and NOT auto-migrated by mall-service.
type SysUser struct {
	ID             int64  `gorm:"primaryKey" json:"id"`
	Username       string `gorm:"type:varchar(64)" json:"username"`
	Password       string `gorm:"type:varchar(255)" json:"-"`
	Mobile         string `gorm:"type:varchar(20)" json:"mobile"`
	RealName       string `gorm:"column:real_name;type:varchar(64)" json:"real_name"`
	Role           string `gorm:"type:varchar(32)" json:"role"`
	Status         int    `json:"status"`
	FaceRegistered bool   `gorm:"column:face_registered" json:"face_registered"`
	FaceImageURL   string `gorm:"column:face_image_url;type:varchar(512)" json:"face_image_url"`
}

func (SysUser) TableName() string { return "sys_user" }
