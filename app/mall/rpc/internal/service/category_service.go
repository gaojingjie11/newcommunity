package service

import (
	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepo
}

func NewCategoryService(repo *repository.CategoryRepo) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) List() ([]model.ProductCategory, error) {
	return s.repo.List()
}

func (s *CategoryService) GetDetail(id int64) (*model.ProductCategory, error) {
	return s.repo.FindByID(id)
}

func (s *CategoryService) Create(cat *model.ProductCategory) error {
	return s.repo.Create(cat)
}

func (s *CategoryService) Update(cat *model.ProductCategory) error {
	return s.repo.Update(cat)
}

func (s *CategoryService) Delete(id int64) error {
	return s.repo.Delete(id)
}
