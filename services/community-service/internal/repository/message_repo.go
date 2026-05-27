package repository

import (
	"time"

	"smartcommunity-microservices/services/community-service/internal/model"

	"gorm.io/gorm"
)

type MessageRepo struct {
	db *gorm.DB
}

func NewMessageRepo(db *gorm.DB) *MessageRepo {
	return &MessageRepo{db: db}
}

func (r *MessageRepo) List(page, size int) ([]model.CommunityMessage, int64, error) {
	var items []model.CommunityMessage
	var total int64

	q := r.db.Model(&model.CommunityMessage{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch messages ordered by created_at ASC (chat chronological order)
	if err := q.Order("created_at ASC").Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	// Batch-load user info for all messages
	if len(items) > 0 {
		userIDs := make([]int64, 0, len(items))
		for _, m := range items {
			userIDs = append(userIDs, m.UserID)
		}
		var users []struct {
			ID       int64
			Username string
			Avatar   string
		}
		r.db.Table("sys_user").Where("id IN ?", userIDs).Select("id, username, avatar").Scan(&users)
		userMap := make(map[int64]*model.CommunityMessageUser, len(users))
		for _, u := range users {
			userMap[u.ID] = &model.CommunityMessageUser{
				ID:       u.ID,
				Username: u.Username,
				Avatar:   u.Avatar,
			}
		}
		for i := range items {
			if u, ok := userMap[items[i].UserID]; ok {
				items[i].User = u
			}
		}
	}

	return items, total, nil
}

func (r *MessageRepo) Create(userID int64, content string) (*model.CommunityMessage, error) {
	msg := &model.CommunityMessage{
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}
	if err := r.db.Create(msg).Error; err != nil {
		return nil, err
	}

	// Load user info for the response
	var user struct {
		ID       int64
		Username string
		Avatar   string
	}
	r.db.Table("sys_user").Where("id = ?", userID).Select("id, username, avatar").Scan(&user)
	msg.User = &model.CommunityMessageUser{
		ID:       user.ID,
		Username: user.Username,
		Avatar:   user.Avatar,
	}

	return msg, nil
}
