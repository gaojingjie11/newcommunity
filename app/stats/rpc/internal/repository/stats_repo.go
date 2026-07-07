package repository

import (
	"fmt"
	"time"

	"smartcommunity-microservices/app/stats/rpc/internal/model"

	"gorm.io/gorm"
)

type StatsRepo struct {
	db *gorm.DB
}

func NewStatsRepo(db *gorm.DB) *StatsRepo {
	return &StatsRepo{db: db}
}

func (r *StatsRepo) ProductSalesRank(limit int) ([]model.ProductSalesRank, error) {
	var results []model.ProductSalesRank
	err := r.db.Raw(`
		SELECT oi.product_id,
		       COALESCE(p.name, '') AS product_name,
		       SUM(oi.quantity)     AS total_sales,
		       SUM(oi.price * oi.quantity) AS total_amount
		FROM oms_order_item oi
		JOIN oms_order o ON o.id = oi.order_id AND o.status IN (1, 2, 3)
		LEFT JOIN pms_product p ON p.id = oi.product_id
		GROUP BY oi.product_id, p.name
		ORDER BY total_amount DESC
		LIMIT ?
	`, limit).Scan(&results).Error
	return results, err
}

func (r *StatsRepo) ProductViewRank(limit int) ([]model.ProductViewRank, error) {
	var results []model.ProductViewRank
	err := r.db.Raw(`
		SELECT id AS product_id,
		       name AS product_name,
		       view_count,
		       0 AS unique_users
		FROM pms_product
		WHERE status = 1
		ORDER BY view_count DESC
		LIMIT ?
	`, limit).Scan(&results).Error
	return results, err
}

func (r *StatsRepo) CommunityOverview() (model.CommunityOverview, error) {
	var ov model.CommunityOverview

	// 1. SysUser total count
	if tx := r.db.Model(&model.SysUserStats{}).Where("status = 1").Count(&ov.UserCount); tx.Error != nil {
		return ov, tx.Error
	}
	ov.TotalUsers = ov.UserCount

	// 2. Order counts
	if tx := r.db.Model(&model.Order{}).Count(&ov.OrderCount); tx.Error != nil {
		return ov, tx.Error
	}

	// 3. PaidAmount from orders
	if tx := r.db.Model(&model.Order{}).Where("status IN (1, 2, 3)").Select("COALESCE(SUM(total_amount), 0)").Scan(&ov.PaidAmount); tx.Error != nil {
		return ov, tx.Error
	}

	// 4. Workorder counts (grouped by type)
	var workorderCounts []struct {
		Type  string
		Count int64
	}
	if tx := r.db.Table("workorders").Select("type, COUNT(*) as count").Group("type").Scan(&workorderCounts); tx.Error == nil {
		for _, wc := range workorderCounts {
			if wc.Type == "repair" {
				ov.RepairCount = wc.Count
			} else if wc.Type == "complaint" {
				ov.ComplaintCount = wc.Count
			}
		}
	}

	// 5. Property fee counts (grouped by status)
	var feeCounts []struct {
		Status int
		Count  int64
	}
	if tx := r.db.Table("property_fees").Select("status, COUNT(*) as count").Group("status").Scan(&feeCounts); tx.Error == nil {
		for _, fc := range feeCounts {
			ov.FeeCount += fc.Count
			if fc.Status == 1 {
				ov.FeePaidCount = fc.Count
			}
		}
	}

	// 6. Today Orders
	if tx := r.db.Model(&model.Order{}).Where("created_at >= CURRENT_DATE").Count(&ov.TodayOrders); tx.Error != nil {
		return ov, tx.Error
	}

	// 7. Parking Rate (grouped by status)
	var totalSpaces, occupiedSpaces int64
	var parkingCounts []struct {
		Status int
		Count  int64
	}
	if tx := r.db.Model(&model.ParkingSpace{}).Select("status, COUNT(*) as count").Group("status").Scan(&parkingCounts); tx.Error == nil {
		for _, pc := range parkingCounts {
			totalSpaces += pc.Count
			if pc.Status == 1 {
				occupiedSpaces = pc.Count
			}
		}
		if totalSpaces > 0 {
			ov.ParkingRate = fmt.Sprintf("%d%%", occupiedSpaces*100/totalSpaces)
		} else {
			ov.ParkingRate = "0%"
		}
	} else {
		ov.ParkingRate = "0%"
	}

	// 8. Month Income (paid orders + paid property fees in current calendar month)
	var orderCents, feeCents int64
	_ = r.db.Model(&model.Order{}).Where("status IN (1, 2, 3) AND paid_at >= DATE_TRUNC('month', CURRENT_DATE)").Select("COALESCE(SUM(total_amount), 0)").Scan(&orderCents)
	_ = r.db.Model(&model.PropertyFee{}).Where("status = 1 AND paid_at >= DATE_TRUNC('month', CURRENT_DATE)").Select("COALESCE(SUM(amount), 0)").Scan(&feeCents)
	ov.MonthIncome = float64(orderCents+feeCents) / 100.0

	// 9. Workorder Stats (grouped by type: repair vs complaint)
	var workorderStats []struct {
		Type  string
		Count int64
	}
	if tx := r.db.Table("workorders").
		Select("type, COUNT(*) as count").
		Group("type").
		Scan(&workorderStats); tx.Error == nil {
		var mappedStats []model.RepairStat
		for _, ws := range workorderStats {
			name := "其他"
			if ws.Type == "repair" {
				name = "报修"
			} else if ws.Type == "complaint" {
				name = "投诉"
			}
			mappedStats = append(mappedStats, model.RepairStat{
				Name:  name,
				Value: ws.Count,
			})
		}
		ov.RepairStats = mappedStats
	}

	// 10. 7-Day Income Trend (Combined daily sums)
	type DailySum struct {
		Date  time.Time `gorm:"column:date"`
		Total int64     `gorm:"column:total"`
	}

	var orderSums []DailySum
	_ = r.db.Raw(`
		SELECT paid_at::date AS date,
		       COALESCE(SUM(total_amount), 0) AS total
		FROM oms_order
		WHERE status IN (1, 2, 3) AND paid_at >= CURRENT_DATE - INTERVAL '6 days'
		GROUP BY paid_at::date
	`).Scan(&orderSums)

	var feeSums []DailySum
	_ = r.db.Raw(`
		SELECT paid_at::date AS date,
		       COALESCE(SUM(amount), 0) AS total
		FROM property_fees
		WHERE status = 1 AND paid_at >= CURRENT_DATE - INTERVAL '6 days'
		GROUP BY paid_at::date
	`).Scan(&feeSums)

	dailyMap := make(map[string]int64)
	for _, s := range orderSums {
		dateStr := s.Date.Format("01-02")
		dailyMap[dateStr] += s.Total
	}
	for _, s := range feeSums {
		dateStr := s.Date.Format("01-02")
		dailyMap[dateStr] += s.Total
	}

	for i := 6; i >= 0; i-- {
		d := time.Now().AddDate(0, 0, -i)
		dateStr := d.Format("01-02")
		totalCents := dailyMap[dateStr]
		ov.IncomeDates = append(ov.IncomeDates, dateStr)
		ov.IncomeTrend = append(ov.IncomeTrend, float64(totalCents)/100.0)
	}

	// 11. Cost Structure ["物业费", "停车费", "商城消费"]
	var totalFeeCents, totalOrderCents int64
	var activeBindings int64
	_ = r.db.Model(&model.PropertyFee{}).Where("status = 1").Select("COALESCE(SUM(amount), 0)").Scan(&totalFeeCents)
	_ = r.db.Model(&model.Order{}).Where("status IN (1, 2, 3)").Select("COALESCE(SUM(total_amount), 0)").Scan(&totalOrderCents)
	_ = r.db.Model(&model.UserParkingBinding{}).Where("status = 1").Count(&activeBindings)

	parkingFee := float64(activeBindings) * 150.0 // 150元/位/月
	ov.CostStructure = []float64{
		float64(totalFeeCents) / 100.0,
		parkingFee,
		float64(totalOrderCents) / 100.0,
	}

	return ov, nil
}

func (r *StatsRepo) OrderSummary() ([]model.OrderSummary, error) {
	var results []model.OrderSummary
	err := r.db.Raw(`
		SELECT status,
		       COUNT(*)              AS count,
		       COALESCE(SUM(total_amount), 0) AS total_amount
		FROM oms_order
		GROUP BY status
		ORDER BY status
	`).Scan(&results).Error
	return results, err
}

func (r *StatsRepo) OrderTrend(days int) ([]model.OrderTrend, error) {
	var results []model.OrderTrend
	err := r.db.Raw(`
		SELECT created_at::date AS date,
		       COUNT(*)          AS count,
		       COALESCE(SUM(total_amount), 0) AS amount
		FROM oms_order
		WHERE created_at >= CURRENT_DATE - (? * INTERVAL '1 day')
		GROUP BY created_at::date
		ORDER BY date
	`, days).Scan(&results).Error
	return results, err
}

func (r *StatsRepo) WorkorderSummary() ([]model.WorkorderSummary, error) {
	var results []model.WorkorderSummary
	err := r.db.Raw(`
		SELECT type, status, COUNT(*) AS count
		FROM workorders
		GROUP BY type, status
	`).Scan(&results).Error
	return results, err
}

func (r *StatsRepo) GreenPointsLeaderboard(limit int) ([]model.EcoLeaderboard, error) {
	var results []model.EcoLeaderboard
	err := r.db.Model(&model.SysUserStats{}).
		Where("status = 1 AND green_points_total_earned > 0").
		Order("green_points_total_earned DESC, id ASC").
		Limit(limit).
		Select("id AS user_id, username, real_name, green_points_total_earned AS green_points").
		Scan(&results).Error
	return results, err
}

func (r *StatsRepo) TotalGreenPointsIssued() (int64, error) {
	var total int64
	err := r.db.Model(&model.SysUserStats{}).
		Where("status = 1").
		Select("COALESCE(SUM(green_points_total_earned), 0)").
		Scan(&total).Error
	return total, err
}
