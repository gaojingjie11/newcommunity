package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"smartcommunity-microservices/services/community-service/internal/model"
	"smartcommunity-microservices/services/community-service/internal/repository"
)

type WorkorderService struct {
	repo *repository.WorkorderRepo
	bus  *EventBus
}

type CreateWorkorderRequest struct {
	Type        string `json:"type" binding:"required"` // repair or complaint
	Category    string `json:"category" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type ProcessRequest struct {
	Status int    `json:"status" binding:"required"`
	Result string `json:"result" binding:"required"`
}

type CreatedResult struct {
	Item        *model.WorkOrder `json:"item"`
	Event       string           `json:"event"`
	EventStatus string           `json:"event_status"`
}

func NewWorkorderService(repo *repository.WorkorderRepo, bus *EventBus) *WorkorderService {
	return &WorkorderService{repo: repo, bus: bus}
}

func (s *WorkorderService) Create(ctx context.Context, userID int64, req CreateWorkorderRequest) (*CreatedResult, error) {
	workorderType := normalizeWorkorderType(req.Type)
	if workorderType == "" {
		return nil, errors.New("type must be repair or complaint")
	}
	item := &model.WorkOrder{
		Type:        workorderType,
		UserID:      userID,
		Category:    strings.TrimSpace(req.Category),
		Description: strings.TrimSpace(req.Description),
		Status:      model.StatusPending,
	}
	if item.Category == "" || item.Description == "" {
		return nil, errors.New("category and description required")
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	event := workorderType + ".created"
	eventStatus := s.bus.Publish(ctx, event, map[string]interface{}{
		"id":          item.ID,
		"type":        item.Type,
		"user_id":     item.UserID,
		"category":    item.Category,
		"description": item.Description,
		"created_at":  item.CreatedAt.Format(time.RFC3339),
	})
	return &CreatedResult{Item: item, Event: event, EventStatus: eventStatus}, nil
}

func (s *WorkorderService) MyList(userID int64, page, size int) ([]model.WorkOrder, int64, error) {
	return s.repo.ListByUser(userID, page, size)
}

func (s *WorkorderService) AdminList(workorderType string, status *int, page, size int) ([]model.WorkOrderAdminView, int64, error) {
	workorderType = normalizeWorkorderType(workorderType)
	return s.repo.ListAll(workorderType, status, page, size)
}

func (s *WorkorderService) Process(id, operatorID int64, req ProcessRequest) (*model.WorkOrder, error) {
	if !validStatus(req.Status) {
		return nil, errors.New("status must be 1 or 2")
	}
	return s.repo.Process(id, operatorID, req.Status, req.Result)
}

func (s *WorkorderService) Logs(targetID int64) ([]model.WorkorderLog, error) {
	return s.repo.ListLogs(targetID)
}

func normalizeWorkorderType(input string) string {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case model.WorkorderTypeRepair, "1", "报修", "设备报修":
		return model.WorkorderTypeRepair
	case model.WorkorderTypeComplaint, "2", "投诉", "建议投诉":
		return model.WorkorderTypeComplaint
	default:
		return ""
	}
}

func validStatus(status int) bool {
	return status == model.StatusProcessing || status == model.StatusCompleted
}
