package service

import (
	"errors"

	"smartcommunity-microservices/services/mall-service/internal/model"
	"smartcommunity-microservices/services/mall-service/internal/repository"

	"gorm.io/gorm"
)

type CartService struct {
	cartRepo    *repository.CartRepo
	productRepo *repository.ProductRepo
}

func NewCartService(cartRepo *repository.CartRepo, productRepo *repository.ProductRepo) *CartService {
	return &CartService{cartRepo: cartRepo, productRepo: productRepo}
}

func (s *CartService) Add(userID, productID int64, quantity int) error {
	if quantity <= 0 {
		return errors.New("数量必须大于0")
	}

	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return errors.New("商品不存在")
	}
	if product.Status != 1 {
		return errors.New("商品已下架")
	}

	existing, err := s.cartRepo.FindByUserProduct(userID, productID)
	if err == nil {
		return s.cartRepo.UpdateQuantity(existing.ID, existing.Quantity+quantity)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return s.cartRepo.Create(&model.Cart{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	})
}

func (s *CartService) Remove(cartID, userID int64) error {
	return s.cartRepo.Delete(cartID, userID)
}

func (s *CartService) UpdateQuantity(cartID, userID, quantity int64) error {
	if quantity <= 0 {
		return s.cartRepo.Delete(cartID, userID)
	}
	return s.cartRepo.UpdateQuantity(cartID, int(quantity))
}

func (s *CartService) List(userID int64) ([]model.Cart, error) {
	return s.cartRepo.ListByUser(userID)
}
