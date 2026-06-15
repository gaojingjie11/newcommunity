package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/repository"

	goredis "github.com/redis/go-redis/v9"
)

type CategoryService struct {
	repo *repository.CategoryRepo
	rdb  *goredis.Client
}

func NewCategoryService(repo *repository.CategoryRepo, rdb *goredis.Client) *CategoryService {
	return &CategoryService{
		repo: repo,
		rdb:  rdb,
	}
}

func (s *CategoryService) List() ([]model.ProductCategory, error) {
	if s.rdb == nil {
		return s.repo.List()
	}

	ctx := context.Background()
	cacheKey := "mall:categories"

	// 1. Try reading category list cache
	cachedJSON, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil && cachedJSON != "" {
		var list []model.ProductCategory
		if err := json.Unmarshal([]byte(cachedJSON), &list); err == nil {
			return list, nil
		}
	}

	// 2. Cache miss — fetch from DB
	list, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	// 3. Cache the list (e.g. 24 hours TTL since categories rarely change)
	listJSON, err := json.Marshal(list)
	if err == nil {
		_ = s.rdb.Set(ctx, cacheKey, string(listJSON), 24*time.Hour).Err()
	}

	return list, nil
}

func (s *CategoryService) GetDetail(id int64) (*model.ProductCategory, error) {
	if s.rdb == nil {
		return s.repo.FindByID(id)
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("mall:category:detail:%d", id)

	// 1. Try reading category detail cache
	cachedJSON, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil && cachedJSON != "" {
		var cat model.ProductCategory
		if err := json.Unmarshal([]byte(cachedJSON), &cat); err == nil {
			return &cat, nil
		}
	}

	// 2. Cache miss — fetch from DB
	cat, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 3. Cache detail (2 hours TTL)
	catJSON, err := json.Marshal(cat)
	if err == nil {
		_ = s.rdb.Set(ctx, cacheKey, string(catJSON), 2*time.Hour).Err()
	}

	return cat, nil
}

func (s *CategoryService) Create(cat *model.ProductCategory) error {
	// Create in DB
	err := s.repo.Create(cat)
	if err != nil {
		return err
	}

	// Evict category list cache
	if s.rdb != nil {
		ctx := context.Background()
		_ = s.rdb.Del(ctx, "mall:categories").Err()
	}

	return nil
}

func (s *CategoryService) Update(cat *model.ProductCategory) error {
	// Update DB first
	err := s.repo.Update(cat)
	if err != nil {
		return err
	}

	// Evict caches to maintain consistency
	if s.rdb != nil {
		ctx := context.Background()
		_ = s.rdb.Del(ctx, "mall:categories", fmt.Sprintf("mall:category:detail:%d", cat.ID)).Err()
	}

	return nil
}

func (s *CategoryService) Delete(id int64) error {
	// Delete DB first
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	// Evict caches to maintain consistency
	if s.rdb != nil {
		ctx := context.Background()
		_ = s.rdb.Del(ctx, "mall:categories", fmt.Sprintf("mall:category:detail:%d", id)).Err()
	}

	return nil
}
