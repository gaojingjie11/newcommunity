package model

import "time"

// ProductComment stores user reviews for mall products.
type ProductComment struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `gorm:"index;not null" json:"user_id"`
	ProductID int64     `gorm:"index;not null" json:"product_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Rating    int       `gorm:"not null;default:5" json:"rating"`
	CreatedAt time.Time `json:"created_at"`

	User UserProfile `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (ProductComment) TableName() string { return "pms_product_comment" }

// UserProfile is a read-only projection of sys_user for comment display.
type UserProfile struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	Username string `gorm:"type:varchar(64)" json:"username"`
	RealName string `gorm:"column:real_name;type:varchar(64)" json:"real_name"`
	Avatar   string `gorm:"type:varchar(255)" json:"avatar"`
}

func (UserProfile) TableName() string { return "sys_user" }
