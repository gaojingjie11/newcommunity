package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/repository"

	goredis "github.com/redis/go-redis/v9"
)

type ProductService struct {
	productRepo *repository.ProductRepo
	rdb         *goredis.Client
	eventBus    *EventBus
}

type ProductListCache struct {
	Products []model.Product `json:"products"`
	Total    int64           `json:"total"`
}

func NewProductService(productRepo *repository.ProductRepo, rdb *goredis.Client, eventBus *EventBus) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		rdb:         rdb,
		eventBus:    eventBus,
	}
}

func (s *ProductService) List(page, size int, categoryID int64, sort string, keyword string) ([]model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	} else if size > 1000 {
		size = 1000
	}

	if s.rdb == nil {
		return s.productRepo.List(page, size, categoryID, sort, keyword)
	}

	ctx := context.Background()
	// Generate distinct key for list cache based on query parameters
	cacheKey := fmt.Sprintf("mall:product:list:cat:%d:page:%d:size:%d:sort:%s:kw:%s", categoryID, page, size, sort, keyword)

	// 1. Try reading list from cache
	cachedJSON, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil && cachedJSON != "" {
		var listCache ProductListCache
		if err := json.Unmarshal([]byte(cachedJSON), &listCache); err == nil {
			return listCache.Products, listCache.Total, nil
		}
	}

	// 2. Cache miss — query PostgreSQL
	products, total, err := s.productRepo.List(page, size, categoryID, sort, keyword)
	if err != nil {
		return nil, 0, err
	}

	// 3. Cache list result for 1 hour
	listCache := ProductListCache{
		Products: products,
		Total:    total,
	}
	listJSON, err := json.Marshal(listCache)
	if err == nil {
		_ = s.rdb.Set(ctx, cacheKey, string(listJSON), 3600*time.Second).Err()
	}

	return products, total, nil
}

func (s *ProductService) AdminList(page, size int, name string, categoryID int64, isPromotion *bool) ([]model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	} else if size > 1000 {
		size = 1000
	}
	return s.productRepo.AdminList(page, size, name, categoryID, isPromotion)
}

func (s *ProductService) Search(keyword string, page, size int) ([]model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	} else if size > 1000 {
		size = 1000
	}
	return s.productRepo.Search(keyword, page, size)
}

func (s *ProductService) GetPromotions(page, size int) ([]model.Product, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	} else if size > 1000 {
		size = 1000
	}
	return s.productRepo.ListPromotionProducts(page, size)
}

func (s *ProductService) GetDetail(id int64) (*model.Product, error) {
	if s.rdb == nil {
		return s.productRepo.FindByID(id)
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("mall:product:detail:%d", id)

	// 1. Try reading from cache
	cachedJSON, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil && cachedJSON != "" {
		var product model.Product
		if err := json.Unmarshal([]byte(cachedJSON), &product); err == nil {
			return &product, nil
		}
	}

	// 2. Cache miss — fetch from PostgreSQL
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 3. Write back to Redis with a randomized TTL (2 hours + up to 10 minutes random offset)
	ttl := 7200 + rand.Intn(600)
	productJSON, err := json.Marshal(product)
	if err == nil {
		_ = s.rdb.Set(ctx, cacheKey, string(productJSON), time.Duration(ttl)*time.Second).Err()
	}

	return product, nil
}

func (s *ProductService) clearListCaches() {
	if s.rdb == nil {
		return
	}
	ctx := context.Background()
	// Safely clear all paginated/sorted list caches in Redis on product changes
	keys, err := s.rdb.Keys(ctx, "mall:product:list:*").Result()
	if err == nil && len(keys) > 0 {
		_ = s.rdb.Del(ctx, keys...).Err()
	}
}

func (s *ProductService) Create(product *model.Product) error {
	err := s.productRepo.Create(product)
	if err != nil {
		return err
	}
	s.clearListCaches()
	return nil
}

func (s *ProductService) Update(product *model.Product) error {
	// Query old product first to get the old image URL
	oldProduct, err := s.productRepo.FindByID(product.ID)
	var oldImageURL string
	if err == nil && oldProduct != nil {
		oldImageURL = oldProduct.ImageURL
	}

	// Update DB first
	err = s.productRepo.Update(product)
	if err != nil {
		return err
	}

	// Delete cache to ensure consistency
	if s.rdb != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("mall:product:detail:%d", product.ID)
		_ = s.rdb.Del(ctx, cacheKey).Err()
	}
	s.clearListCaches()

	// If image URL has changed, publish cleanup event for the old image
	if oldImageURL != "" && oldImageURL != product.ImageURL && s.eventBus != nil {
		s.eventBus.PublishFileCleanup(oldImageURL)
	}

	return nil
}

func (s *ProductService) UpdateFields(id int64, fields map[string]interface{}) error {
	// Update DB first
	err := s.productRepo.UpdateFields(id, fields)
	if err != nil {
		return err
	}

	// Delete cache to ensure consistency
	if s.rdb != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("mall:product:detail:%d", id)
		_ = s.rdb.Del(ctx, cacheKey).Err()
	}
	s.clearListCaches()

	return nil
}

func (s *ProductService) Delete(id int64) error {
	// Query old product first to get the image URL
	oldProduct, err := s.productRepo.FindByID(id)
	if err != nil {
		// If product doesn't exist, let delete logic handle it
		return s.productRepo.Delete(id)
	}

	err = s.productRepo.Delete(id)
	if err != nil {
		return err
	}

	if s.rdb != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("mall:product:detail:%d", id)
		_ = s.rdb.Del(ctx, cacheKey).Err()
	}
	s.clearListCaches()

	// Publish file cleanup if image exists
	if oldProduct.ImageURL != "" && s.eventBus != nil {
		s.eventBus.PublishFileCleanup(oldProduct.ImageURL)
	}

	return nil
}
