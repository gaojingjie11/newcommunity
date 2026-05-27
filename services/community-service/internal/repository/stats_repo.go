package repository

import (
	"smartcommunity-microservices/services/community-service/internal/model"

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
		SELECT vl.product_id,
		       COALESCE(p.name, '') AS product_name,
		       COUNT(*)             AS view_count,
		       COUNT(DISTINCT CASE WHEN vl.user_id > 0 THEN vl.user_id ELSE NULL END) AS unique_users
		FROM product_view_logs vl
		LEFT JOIN pms_product p ON p.id = vl.product_id
		GROUP BY vl.product_id, p.name
		ORDER BY view_count DESC
		LIMIT ?
	`, limit).Scan(&results).Error
	return results, err
}

func (r *StatsRepo) CommunityOverview() (model.CommunityOverview, error) {
	var ov model.CommunityOverview

	if tx := r.db.Model(&model.SysUserStats{}).Where("status = 1").Count(&ov.UserCount); tx.Error != nil {
		return ov, tx.Error
	}

	if tx := r.db.Model(&model.Order{}).Count(&ov.OrderCount); tx.Error != nil {
		return ov, tx.Error
	}

	if tx := r.db.Model(&model.Order{}).Where("status IN (1, 2, 3)").Select("COALESCE(SUM(total_amount), 0)").Scan(&ov.PaidAmount); tx.Error != nil {
		return ov, tx.Error
	}

	if tx := r.db.Model(&model.WorkOrder{}).Where("type = ?", model.WorkorderTypeRepair).Count(&ov.RepairCount); tx.Error != nil {
		return ov, tx.Error
	}

	if tx := r.db.Model(&model.WorkOrder{}).Where("type = ?", model.WorkorderTypeComplaint).Count(&ov.ComplaintCount); tx.Error != nil {
		return ov, tx.Error
	}

	if tx := r.db.Model(&model.PropertyFee{}).Count(&ov.FeeCount); tx.Error != nil {
		return ov, tx.Error
	}
	if tx := r.db.Model(&model.PropertyFee{}).Where("status = 1").Count(&ov.FeePaidCount); tx.Error != nil {
		return ov, tx.Error
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
		SELECT DATE(created_at) AS date,
		       COUNT(*)          AS count,
		       COALESCE(SUM(total_amount), 0) AS amount
		FROM oms_order
		WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
		GROUP BY DATE(created_at)
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
