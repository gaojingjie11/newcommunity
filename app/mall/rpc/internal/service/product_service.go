package service

import (
	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/repository"
)

type ProductService struct {
	productRepo *repository.ProductRepo
}

func NewProductService(productRepo *repository.ProductRepo) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) List(page, size int, categoryID int64, sort string, keyword string) ([]model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return s.productRepo.List(page, size, categoryID, sort, keyword)
}

func (s *ProductService) AdminList(page, size int, name string, categoryID int64, isPromotion *bool) ([]model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return s.productRepo.AdminList(page, size, name, categoryID, isPromotion)
}

func (s *ProductService) Search(keyword string, page, size int) ([]model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return s.productRepo.Search(keyword, page, size)
}

func (s *ProductService) GetPromotions(page, size int) ([]model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return s.productRepo.ListPromotionProducts(page, size)
}

func (s *ProductService) GetDetail(id int64) (*model.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *ProductService) Create(product *model.Product) error {
	return s.productRepo.Create(product)
}

func (s *ProductService) Update(product *model.Product) error {
	return s.productRepo.Update(product)
}

func (s *ProductService) UpdateFields(id int64, fields map[string]interface{}) error {
	return s.productRepo.UpdateFields(id, fields)
}

func (s *ProductService) Delete(id int64) error {
	return s.productRepo.Delete(id)
}
