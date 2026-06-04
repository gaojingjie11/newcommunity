package service

import (
	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/repository"
)

type StoreService struct {
	storeRepo *repository.StoreRepo
	spRepo    *repository.StoreProductRepo
}

func NewStoreService(storeRepo *repository.StoreRepo, spRepo *repository.StoreProductRepo) *StoreService {
	return &StoreService{storeRepo: storeRepo, spRepo: spRepo}
}

func (s *StoreService) List(page, size int, areaID int64) ([]model.Store, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return s.storeRepo.List(page, size, areaID)
}

func (s *StoreService) GetDetail(id int64) (*model.Store, error) {
	return s.storeRepo.FindByID(id)
}

func (s *StoreService) Create(store *model.Store) error {
	return s.storeRepo.Create(store)
}

func (s *StoreService) Update(store *model.Store) error {
	return s.storeRepo.Update(store)
}

func (s *StoreService) Delete(id int64) error {
	return s.storeRepo.Delete(id)
}

func (s *StoreService) BindProduct(storeID, productID int64, stock int) error {
	return s.spRepo.Bind(storeID, productID, stock)
}

func (s *StoreService) UnbindProduct(storeID, productID int64) error {
	return s.spRepo.Unbind(storeID, productID)
}

func (s *StoreService) UpdateProductStatus(storeID, productID int64, status int) error {
	return s.spRepo.UpdateStatus(storeID, productID, status)
}

func (s *StoreService) UpdateProductStock(storeID, productID int64, stock int) error {
	return s.spRepo.UpdateStock(storeID, productID, stock)
}

func (s *StoreService) ListProducts(storeID int64) ([]model.StoreProduct, error) {
	return s.spRepo.ListByStore(storeID)
}
