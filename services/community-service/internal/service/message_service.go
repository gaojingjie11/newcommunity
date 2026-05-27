package service

import (
	"errors"

	"smartcommunity-microservices/services/community-service/internal/model"
	"smartcommunity-microservices/services/community-service/internal/repository"
)

type MessageService struct {
	repo *repository.MessageRepo
}

func NewMessageService(repo *repository.MessageRepo) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) ListMessages(page, size int) ([]model.CommunityMessage, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 50
	}
	return s.repo.List(page, size)
}

func (s *MessageService) SendMessage(userID int64, content string) (*model.CommunityMessage, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}
	if content == "" {
		return nil, errors.New("content is required")
	}
	return s.repo.Create(userID, content)
}
