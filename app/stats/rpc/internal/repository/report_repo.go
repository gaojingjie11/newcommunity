package repository

import (
	"smartcommunity-microservices/app/stats/rpc/internal/model"

	"gorm.io/gorm"
)

type ReportRepo struct {
	db *gorm.DB
}

func NewReportRepo(db *gorm.DB) *ReportRepo {
	return &ReportRepo{db: db}
}

func (r *ReportRepo) Create(report *model.AIReport) error {
	return r.db.Create(report).Error
}

func (r *ReportRepo) FindLatest() (*model.AIReport, error) {
	var report model.AIReport
	err := r.db.Order("id DESC").First(&report).Error
	return &report, err
}

func (r *ReportRepo) FindByID(id int64) (*model.AIReport, error) {
	var report model.AIReport
	err := r.db.Where("id = ?", id).First(&report).Error
	return &report, err
}

func (r *ReportRepo) List(page, size int) ([]model.AIReport, int64, error) {
	var reports []model.AIReport
	var total int64

	query := r.db.Model(&model.AIReport{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("id DESC").Find(&reports).Error; err != nil {
		return nil, 0, err
	}
	return reports, total, nil
}

func (r *ReportRepo) Count7DayRepairs() (int64, error) {
	var count int64
	err := r.db.Table("workorders").
		Where("type = ? AND created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)", "repair").
		Count(&count).Error
	return count, err
}

func (r *ReportRepo) CountPendingRepairs() (int64, error) {
	var count int64
	err := r.db.Table("workorders").
		Where("type = ? AND status <> ?", "repair", 2). // 2 = completed
		Count(&count).Error
	return count, err
}

func (r *ReportRepo) Count7DayVisitors() (int64, error) {
	var count int64
	err := r.db.Table("visitors").Where("created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)").Count(&count).Error
	return count, err
}

func (r *ReportRepo) Count7DayPaidFees() (int64, error) {
	var count int64
	err := r.db.Table("property_fees").Where("status = 1 AND paid_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)").Count(&count).Error
	return count, err
}

func (r *ReportRepo) Sum7DayPaidAmount() (float64, error) {
	var totalCents int64
	err := r.db.Table("property_fees").
		Where("status = 1 AND paid_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalCents).Error
	if err != nil {
		return 0, err
	}
	return float64(totalCents) / 100.0, nil
}
