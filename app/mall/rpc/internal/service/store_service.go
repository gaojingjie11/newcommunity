package service

import (
	"fmt"
	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/repository"
)

type StoreService struct {
	storeRepo   *repository.StoreRepo
	spRepo      *repository.StoreProductRepo
	productRepo *repository.ProductRepo
}

func NewStoreService(storeRepo *repository.StoreRepo, spRepo *repository.StoreProductRepo, productRepo *repository.ProductRepo) *StoreService {
	return &StoreService{storeRepo: storeRepo, spRepo: spRepo, productRepo: productRepo}
}

func (s *StoreService) List(page, size int, areaID int64, storeIDs []int64) ([]model.Store, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return s.storeRepo.List(page, size, areaID, storeIDs)
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
	if err := s.validateStoreProductStock(storeID, productID, stock); err != nil {
		return err
	}
	return s.spRepo.Bind(storeID, productID, stock)
}

func (s *StoreService) UnbindProduct(storeID, productID int64) error {
	return s.spRepo.Unbind(storeID, productID)
}

func (s *StoreService) UpdateProductStatus(storeID, productID int64, status int) error {
	return s.spRepo.UpdateStatus(storeID, productID, status)
}

func (s *StoreService) UpdateProductStock(storeID, productID int64, stock int) error {
	if err := s.validateStoreProductStock(storeID, productID, stock); err != nil {
		return err
	}
	return s.spRepo.UpdateStock(storeID, productID, stock)
}

func (s *StoreService) ListProducts(storeID int64) ([]model.StoreProduct, error) {
	return s.spRepo.ListByStore(storeID)
}

func (s *StoreService) validateStoreProductStock(storeID, productID int64, stock int) error {
	if stock < 0 {
		return fmt.Errorf("分配库存不能小于 0")
	}

	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return err
	}

	if stock > product.Stock {
		return fmt.Errorf("分配库存不能超过商品总库存（当前总库存: %d）", product.Stock)
	}

	allocatedOtherStores, err := s.spRepo.SumAllocatedStock(productID, storeID)
	if err != nil {
		return err
	}

	if allocatedOtherStores+stock > product.Stock {
		remaining := product.Stock - allocatedOtherStores
		if remaining < 0 {
			remaining = 0
		}
		return fmt.Errorf("库存分配超出商品总库存，当前其他门店已分配 %d，当前门店最多还能分配 %d", allocatedOtherStores, remaining)
	}

	return nil
}
