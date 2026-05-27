package service

import (
	"errors"
	"strings"

	"smartcommunity-microservices/services/community-service/internal/model"
	"smartcommunity-microservices/services/community-service/internal/repository"
)

type ParkingService struct {
	repo *repository.ParkingRepo
}

type CreateParkingRequest struct {
	ParkingNo string `json:"parking_no" binding:"required"`
}

type AssignParkingRequest struct {
	Mobile   string `json:"mobile"`
	UserID   int64  `json:"user_id,omitempty"`
	CarPlate string `json:"car_plate"`
}

type BindPlateRequest struct {
	CarPlate string `json:"car_plate" binding:"required"`
}

func NewParkingService(repo *repository.ParkingRepo) *ParkingService {
	return &ParkingService{repo: repo}
}

func (s *ParkingService) List(page, size int) ([]model.ParkingSpaceAdminView, int64, error) {
	return s.repo.List(page, size)
}

func (s *ParkingService) Create(req CreateParkingRequest) (*model.ParkingSpace, error) {
	parkingNo := strings.TrimSpace(req.ParkingNo)
	if parkingNo == "" {
		return nil, errors.New("parking_no required")
	}
	item := &model.ParkingSpace{ParkingNo: parkingNo, Status: 0}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *ParkingService) MyBindings(userID int64) ([]model.UserParkingBinding, error) {
	return s.repo.ListBindingsByUser(userID)
}

func (s *ParkingService) Assign(spaceID int64, req AssignParkingRequest) (*model.UserParkingBinding, error) {
	if strings.TrimSpace(req.Mobile) == "" && req.UserID > 0 {
		return nil, errors.New("mobile required")
	}
	return s.repo.Assign(spaceID, req.Mobile, req.CarPlate)
}

func (s *ParkingService) BindPlate(bindingID, userID int64, req BindPlateRequest) (*model.UserParkingBinding, error) {
	return s.repo.BindPlate(bindingID, userID, req.CarPlate)
}

func (s *ParkingService) Stats() (map[string]int64, error) {
	return s.repo.Stats()
}
