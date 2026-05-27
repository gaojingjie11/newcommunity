package service

import (
	"encoding/json"
	"log/slog"
	"os"
	"testing"
	"time"

	"smartcommunity-microservices/services/community-service/internal/model"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

// ── Cache helper tests ──

func TestGetJSONCache_Miss(t *testing.T) {
	s := miniredis.RunT(t)
	rdb := goredis.NewClient(&goredis.Options{Addr: s.Addr()})
	defer rdb.Close()

	svc := &StatsService{rdb: rdb, log: testLogger()}

	var dest []model.ProductSalesRank
	if svc.getJSONCache(t.Context(), "nonexistent:key", &dest) {
		t.Fatal("expected cache miss")
	}
}

func TestSetAndGetJSONCache_Hit(t *testing.T) {
	s := miniredis.RunT(t)
	rdb := goredis.NewClient(&goredis.Options{Addr: s.Addr()})
	defer rdb.Close()

	svc := &StatsService{rdb: rdb, log: testLogger()}

	original := []model.ProductSalesRank{
		{ProductID: 1, ProductName: "Apple", TotalSales: 10, TotalAmount: 1000},
		{ProductID: 2, ProductName: "Banana", TotalSales: 5, TotalAmount: 500},
	}
	svc.setJSONCache(t.Context(), "test:key", original, 30*time.Second)

	var cached []model.ProductSalesRank
	if !svc.getJSONCache(t.Context(), "test:key", &cached) {
		t.Fatal("expected cache hit")
	}
	if len(cached) != 2 || cached[0].ProductName != "Apple" {
		t.Fatalf("unexpected cached data: %+v", cached)
	}
}

func TestGetJSONCache_InvalidJSON(t *testing.T) {
	s := miniredis.RunT(t)
	rdb := goredis.NewClient(&goredis.Options{Addr: s.Addr()})
	defer rdb.Close()

	// Manually set invalid JSON
	s.Set("bad:key", "not-json")

	svc := &StatsService{rdb: rdb, log: testLogger()}
	var dest []model.ProductSalesRank
	if svc.getJSONCache(t.Context(), "bad:key", &dest) {
		t.Fatal("expected cache miss due to invalid JSON")
	}
}

func TestCacheHelpers_NilRedis(t *testing.T) {
	svc := &StatsService{rdb: nil, log: testLogger()}

	// Should not panic
	var dest []model.ProductSalesRank
	if svc.getJSONCache(t.Context(), "any:key", &dest) {
		t.Fatal("expected false with nil redis")
	}
	svc.setJSONCache(t.Context(), "any:key", dest, 30*time.Second) // should be no-op
}

func TestCacheExpiry(t *testing.T) {
	s := miniredis.RunT(t)
	rdb := goredis.NewClient(&goredis.Options{Addr: s.Addr()})
	defer rdb.Close()

	svc := &StatsService{rdb: rdb, log: testLogger()}

	data := model.CommunityOverview{UserCount: 100}
	svc.setJSONCache(t.Context(), "test:expire", data, 1*time.Second)

	var cached model.CommunityOverview
	if !svc.getJSONCache(t.Context(), "test:expire", &cached) {
		t.Fatal("expected cache hit before expiry")
	}

	// Fast-forward miniredis time
	s.FastForward(2 * time.Second)

	if svc.getJSONCache(t.Context(), "test:expire", &cached) {
		t.Fatal("expected cache miss after expiry")
	}
}

// ── JSON serialization round-trip tests ──

func TestJSONRoundTrip_ProductSalesRank(t *testing.T) {
	original := []model.ProductSalesRank{
		{ProductID: 1, ProductName: "Test", TotalSales: 100, TotalAmount: 9999},
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	var decoded []model.ProductSalesRank
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded[0].ProductID != 1 || decoded[0].TotalAmount != 9999 {
		t.Fatalf("round-trip mismatch: %+v", decoded)
	}
}

func TestJSONRoundTrip_CommunityOverview(t *testing.T) {
	original := model.CommunityOverview{
		UserCount: 50, OrderCount: 200, PaidAmount: 50000,
		RepairCount: 10, ComplaintCount: 3, FeeCount: 100, FeePaidCount: 80,
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	var decoded model.CommunityOverview
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.UserCount != 50 || decoded.FeePaidCount != 80 {
		t.Fatalf("round-trip mismatch: %+v", decoded)
	}
}

func TestJSONRoundTrip_OrderStatsCombined(t *testing.T) {
	type orderCache struct {
		Summary []model.OrderSummary `json:"summary"`
		Trend   []model.OrderTrend   `json:"trend"`
	}
	original := orderCache{
		Summary: []model.OrderSummary{{Status: 1, Count: 10, TotalAmount: 5000}},
		Trend:   []model.OrderTrend{{Date: "2025-01-01", Count: 5, Amount: 2500}},
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	var decoded orderCache
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.Summary) != 1 || len(decoded.Trend) != 1 {
		t.Fatalf("round-trip mismatch: %+v", decoded)
	}
}

func TestJSONRoundTrip_WorkorderSummary(t *testing.T) {
	original := []model.WorkorderSummary{
		{Type: "repair", Status: 0, Count: 5},
		{Type: "complaint", Status: 1, Count: 3},
	}
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	var decoded []model.WorkorderSummary
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded[0].Type != "repair" || decoded[1].Count != 3 {
		t.Fatalf("round-trip mismatch: %+v", decoded)
	}
}
