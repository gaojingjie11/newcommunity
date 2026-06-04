package service

import (
	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/repository"
)

type PromotionService struct {
	promoRepo *repository.PromotionRepo
}

func NewPromotionService(promoRepo *repository.PromotionRepo) *PromotionService {
	return &PromotionService{promoRepo: promoRepo}
}

func (s *PromotionService) List(page, size int) ([]model.Promotion, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return s.promoRepo.List(page, size)
}

func (s *PromotionService) GetDetail(id int64) (*model.Promotion, error) {
	return s.promoRepo.FindByID(id)
}

func (s *PromotionService) Create(promo *model.Promotion) error {
	return s.promoRepo.Create(promo)
}

func (s *PromotionService) Update(promo *model.Promotion) error {
	return s.promoRepo.Update(promo)
}

func (s *PromotionService) Delete(id int64) error {
	return s.promoRepo.Delete(id)
}

func (s *PromotionService) BindProducts(promotionID int64, productIDs []int64) error {
	return s.promoRepo.BindProducts(promotionID, productIDs)
}
