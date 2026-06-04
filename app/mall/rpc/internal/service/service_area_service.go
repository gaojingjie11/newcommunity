package service

import (
	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/repository"
)

type ServiceAreaService struct {
	repo *repository.ServiceAreaRepo
}

func NewServiceAreaService(repo *repository.ServiceAreaRepo) *ServiceAreaService {
	return &ServiceAreaService{repo: repo}
}

func (s *ServiceAreaService) List() ([]model.ServiceArea, error) {
	return s.repo.List()
}

func (s *ServiceAreaService) ListAll() ([]model.ServiceArea, error) {
	return s.repo.ListAll()
}

func (s *ServiceAreaService) GetDetail(id int64) (*model.ServiceArea, error) {
	return s.repo.FindByID(id)
}

func (s *ServiceAreaService) Create(area *model.ServiceArea) error {
	return s.repo.Create(area)
}

func (s *ServiceAreaService) Update(area *model.ServiceArea) error {
	return s.repo.Update(area)
}

func (s *ServiceAreaService) Delete(id int64) error {
	return s.repo.Delete(id)
}
