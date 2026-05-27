package service

import (
	"errors"
	"time"

	"smartcommunity-microservices/services/community-service/internal/model"
	"smartcommunity-microservices/services/community-service/internal/repository"
)

type VisitorService struct {
	repo *repository.VisitorRepo
}

type CreateVisitorRequest struct {
	VisitorName  string `json:"visitor_name" binding:"required"`
	VisitorPhone string `json:"visitor_phone" binding:"required"`
	VisitPurpose string `json:"visit_purpose" binding:"required"`
	ReleaseTime  string `json:"release_time" binding:"required"`
	ValidDate    string `json:"valid_date" binding:"required"`
}

type AuditVisitorRequest struct {
	Status int    `json:"status" binding:"required"`
	Remark string `json:"remark"`
}

func NewVisitorService(repo *repository.VisitorRepo) *VisitorService {
	return &VisitorService{repo: repo}
}

func (s *VisitorService) Create(userID int64, req CreateVisitorRequest) (*model.Visitor, error) {
	releaseTime, err := parseTime(req.ReleaseTime)
	if err != nil {
		return nil, err
	}
	validDate, err := parseDate(req.ValidDate)
	if err != nil {
		return nil, err
	}
	item := &model.Visitor{
		UserID:       userID,
		VisitorName:  req.VisitorName,
		VisitorPhone: req.VisitorPhone,
		VisitPurpose: req.VisitPurpose,
		ReleaseTime:  releaseTime,
		ValidDate:    validDate,
		Status:       0,
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *VisitorService) ListByUser(userID int64, page, size int) ([]model.Visitor, int64, error) {
	return s.repo.ListByUser(userID, page, size)
}

func (s *VisitorService) ListAll(status *int, page, size int) ([]model.VisitorAdminView, int64, error) {
	return s.repo.ListAll(status, page, size)
}

func (s *VisitorService) Audit(id int64, req AuditVisitorRequest) (*model.Visitor, error) {
	if req.Status != 1 && req.Status != 2 {
		return nil, errors.New("status must be 1 or 2")
	}
	return s.repo.Audit(id, req.Status, req.Remark)
}

func parseTime(input string) (time.Time, error) {
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", input, time.Local); err == nil {
		return t, nil
	}
	return time.ParseInLocation(time.RFC3339, input, time.Local)
}

func parseDate(input string) (time.Time, error) {
	if t, err := time.ParseInLocation("2006-01-02", input, time.Local); err == nil {
		return t, nil
	}
	return parseTime(input)
}
