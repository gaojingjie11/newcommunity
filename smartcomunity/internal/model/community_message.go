package model

import "time"

// CommunityMessage represents a global community chat message visible to all users.
type CommunityMessage struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `gorm:"column:user_id;index;not null" json:"user_id"`
	Content   string    `gorm:"column:content;type:longtext;not null" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`

	User SysUser `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (CommunityMessage) TableName() string {
	return "cms_community_messages"
}
