package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"smartcommunity-microservices/app/stats/rpc/internal/model"
	"smartcommunity-microservices/app/stats/rpc/internal/repository"

	goredis "github.com/redis/go-redis/v9"
)

const (
	cacheKeySalesRank  = "stats:product:sales-rank:%d"
	cacheKeyViewRank   = "stats:product:view-rank:%d"
	cacheKeyOverview   = "stats:community:overview"
	cacheKeyOrders     = "stats:orders:%d"
	cacheKeyWorkorders = "stats:workorders:summary"
	cacheKeyLeaderboard = "stats:green:leaderboard:%d"

	ttlStatsShort = 30 * time.Second
	ttlStatsLong  = 60 * time.Second
)

type StatsService struct {
	repo *repository.StatsRepo
	rdb  *goredis.Client
	log  *slog.Logger
}

func NewStatsService(repo *repository.StatsRepo, rdb *goredis.Client, log *slog.Logger) *StatsService {
	return &StatsService{repo: repo, rdb: rdb, log: log}
}

func (s *StatsService) getJSONCache(ctx context.Context, key string, dest interface{}) bool {
	if s.rdb == nil {
		return false
	}
	data, err := s.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return false
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return false
	}
	return true
}

func (s *StatsService) setJSONCache(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	if s.rdb == nil {
		return
	}
	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	if err := s.rdb.Set(ctx, key, data, ttl).Err(); err != nil && s.log != nil {
		s.log.Warn("redis cache set failed", "key", key, "error", err)
	}
}

func (s *StatsService) ProductSalesRank(limit int) ([]model.ProductSalesRank, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	ctx := context.Background()
	key := fmt.Sprintf(cacheKeySalesRank, limit)

	var cached []model.ProductSalesRank
	if s.getJSONCache(ctx, key, &cached) {
		return cached, nil
	}

	result, err := s.repo.ProductSalesRank(limit)
	if err != nil {
		return nil, err
	}
	s.setJSONCache(ctx, key, result, ttlStatsLong)
	return result, nil
}

func (s *StatsService) ProductViewRank(limit int) ([]model.ProductViewRank, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	ctx := context.Background()
	key := fmt.Sprintf(cacheKeyViewRank, limit)

	var cached []model.ProductViewRank
	if s.getJSONCache(ctx, key, &cached) {
		return cached, nil
	}

	result, err := s.repo.ProductViewRank(limit)
	if err != nil {
		return nil, err
	}
	s.setJSONCache(ctx, key, result, ttlStatsShort)
	return result, nil
}

func (s *StatsService) CommunityOverview() (model.CommunityOverview, error) {
	ctx := context.Background()

	var cached model.CommunityOverview
	if s.getJSONCache(ctx, cacheKeyOverview, &cached) {
		return cached, nil
	}

	result, err := s.repo.CommunityOverview()
	if err != nil {
		return result, err
	}
	s.setJSONCache(ctx, cacheKeyOverview, result, ttlStatsShort)
	return result, nil
}

func (s *StatsService) OrderSummary() ([]model.OrderSummary, error) {
	return s.repo.OrderSummary()
}

func (s *StatsService) OrderTrend(days int) ([]model.OrderTrend, error) {
	if days <= 0 || days > 90 {
		days = 30
	}
	return s.repo.OrderTrend(days)
}

func (s *StatsService) OrderStatsCombined(days int) (summary []model.OrderSummary, trend []model.OrderTrend, err error) {
	if days <= 0 || days > 90 {
		days = 30
	}
	ctx := context.Background()
	key := fmt.Sprintf(cacheKeyOrders, days)

	type orderCache struct {
		Summary []model.OrderSummary `json:"summary"`
		Trend   []model.OrderTrend   `json:"trend"`
	}

	var cached orderCache
	if s.getJSONCache(ctx, key, &cached) {
		return cached.Summary, cached.Trend, nil
	}

	summary, err = s.repo.OrderSummary()
	if err != nil {
		return nil, nil, err
	}
	trend, err = s.repo.OrderTrend(days)
	if err != nil {
		return nil, nil, err
	}

	s.setJSONCache(ctx, key, orderCache{Summary: summary, Trend: trend}, ttlStatsLong)
	return summary, trend, nil
}

func (s *StatsService) WorkorderSummary() ([]model.WorkorderSummary, error) {
	ctx := context.Background()

	var cached []model.WorkorderSummary
	if s.getJSONCache(ctx, cacheKeyWorkorders, &cached) {
		return cached, nil
	}

	result, err := s.repo.WorkorderSummary()
	if err != nil {
		return nil, err
	}
	s.setJSONCache(ctx, cacheKeyWorkorders, result, ttlStatsLong)
	return result, nil
}

func (s *StatsService) GreenPointsLeaderboard(limit int) ([]model.EcoLeaderboard, int64, error) {
	ctx := context.Background()
	key := fmt.Sprintf(cacheKeyLeaderboard, limit)

	type leaderboardCache struct {
		List              []model.EcoLeaderboard `json:"list"`
		TotalPointsIssued int64                  `json:"total_points_issued"`
	}

	var cached leaderboardCache
	if s.getJSONCache(ctx, key, &cached) {
		return cached.List, cached.TotalPointsIssued, nil
	}

	list, err := s.repo.GreenPointsLeaderboard(limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.TotalGreenPointsIssued()
	if err != nil {
		return nil, 0, err
	}

	s.setJSONCache(ctx, key, leaderboardCache{List: list, TotalPointsIssued: total}, ttlStatsLong)
	return list, total, nil
}
