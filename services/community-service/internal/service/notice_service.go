package service

import (
	"errors"

	"smartcommunity-microservices/services/community-service/internal/model"
	"smartcommunity-microservices/services/community-service/internal/repository"

	"gorm.io/gorm"
)

type NoticeService struct {
	repo *repository.NoticeRepo
}

type CreateNoticeRequest struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Publisher string `json:"publisher"`
}

func NewNoticeService(repo *repository.NoticeRepo) *NoticeService {
	return &NoticeService{repo: repo}
}

func (s *NoticeService) List(page, size int, includeHidden bool) ([]model.Notice, int64, error) {
	return s.repo.List(page, size, includeHidden)
}

func (s *NoticeService) Get(id, viewerID int64) (*model.Notice, error) {
	return s.repo.View(id, viewerID)
}

func (s *NoticeService) Create(req CreateNoticeRequest) (*model.Notice, error) {
	if req.Publisher == "" {
		req.Publisher = "admin"
	}
	item := &model.Notice{
		Title:     req.Title,
		Content:   req.Content,
		Publisher: req.Publisher,
		Status:    1,
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *NoticeService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *NoticeService) MarkRead(id, userID int64) error {
	if userID <= 0 {
		return errors.New("invalid user id")
	}
	if _, err := s.repo.Get(id); err != nil {
		return err
	}
	return s.repo.MarkRead(id, userID)
}

func (s *NoticeService) ListViews(noticeID int64, page, size int) ([]model.NoticeViewLog, int64, error) {
	if _, err := s.repo.Get(noticeID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, err
		}
		return nil, 0, err
	}
	return s.repo.ListViews(noticeID, page, size)
}
