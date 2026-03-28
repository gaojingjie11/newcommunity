package service

import (
	"errors"
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
	"strings"

	"gorm.io/gorm"
)

type CommunityMessageService struct{}

func preloadCommunityUser(db *gorm.DB) *gorm.DB {
	return db.Select("id", "username", "avatar")
}

func (s *CommunityMessageService) Send(userID int64, content string) (*model.CommunityMessage, error) {
	text := strings.TrimSpace(content)
	if text == "" {
		return nil, errors.New("message content cannot be empty")
	}
	if len([]rune(text)) > 1000 {
		return nil, errors.New("message is too long")
	}

	msg := &model.CommunityMessage{
		UserID:  userID,
		Content: text,
	}
	if err := global.DB.Create(msg).Error; err != nil {
		return nil, err
	}
	if err := global.DB.Preload("User", preloadCommunityUser).First(msg, msg.ID).Error; err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *CommunityMessageService) List(page, size int) ([]model.CommunityMessage, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	var list []model.CommunityMessage
	var total int64

	db := global.DB.Model(&model.CommunityMessage{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := global.DB.
		Preload("User", preloadCommunityUser).
		Order("created_at desc").
		Offset(offset).
		Limit(size).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
