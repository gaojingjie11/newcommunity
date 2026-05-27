package service

import (
	"errors"

	"smartcommunity-microservices/services/mall-service/internal/model"
	"smartcommunity-microservices/services/mall-service/internal/repository"
)

type FavoriteService struct {
	favRepo *repository.FavoriteRepo
}

func NewFavoriteService(favRepo *repository.FavoriteRepo) *FavoriteService {
	return &FavoriteService{favRepo: favRepo}
}

// MALL-008: Add favorite
func (s *FavoriteService) Add(userID, productID int64) error {
	exists, err := s.favRepo.Exists(userID, productID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("已收藏")
	}
	return s.favRepo.Add(userID, productID)
}

// MALL-009: Remove favorite
func (s *FavoriteService) Remove(userID, productID int64) error {
	return s.favRepo.Remove(userID, productID)
}

// MALL-018: List favorites
func (s *FavoriteService) List(userID int64, page, size int) ([]model.Favorite, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return s.favRepo.ListByUser(userID, page, size)
}

func (s *FavoriteService) Check(userID, productID int64) (bool, error) {
	return s.favRepo.Exists(userID, productID)
}
