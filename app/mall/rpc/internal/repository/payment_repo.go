package repository

import (
	"smartcommunity-microservices/app/mall/rpc/internal/consts"
	"smartcommunity-microservices/app/mall/rpc/internal/model"

	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) WithTx(tx *gorm.DB) *PaymentRepo {
	return &PaymentRepo{db: tx}
}

func (r *PaymentRepo) Create(record *model.PaymentRecord) error {
	return r.db.Create(record).Error
}

func (r *PaymentRepo) FindByIdempotencyKey(key string) (*model.PaymentRecord, error) {
	var record model.PaymentRecord
	err := r.db.Where("idempotency_key = ?", key).First(&record).Error
	return &record, err
}

func (r *PaymentRepo) FindByOrderID(orderID int64) (*model.PaymentRecord, error) {
	var record model.PaymentRecord
	err := r.db.Where("order_id = ?", orderID).Order("id desc").First(&record).Error
	return &record, err
}

func (r *PaymentRepo) FindByOrderNo(orderNo string) (*model.PaymentRecord, error) {
	var record model.PaymentRecord
	err := r.db.Where("order_no = ?", orderNo).First(&record).Error
	return &record, err
}

func (r *PaymentRepo) ListPendingRechargeRecords(limit int) ([]model.PaymentRecord, error) {
	var records []model.PaymentRecord
	q := r.db.
		Where("status = ? AND payment_method = ? AND order_no LIKE ?", consts.PaymentStatusInit, "alipay", "RECH_%").
		Order("id asc")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&records).Error
	return records, err
}

func (r *PaymentRepo) ListPendingRechargeRecordsByUser(userID int64, limit int) ([]model.PaymentRecord, error) {
	var records []model.PaymentRecord
	q := r.db.
		Where("user_id = ? AND status = ? AND payment_method = ? AND order_no LIKE ?", userID, consts.PaymentStatusInit, "alipay", "RECH_%").
		Order("id desc")
	if limit > 0 {
		q = q.Limit(limit)
	}
	err := q.Find(&records).Error
	return records, err
}

// UpdateStatus performs a conditional status update on a payment record.
func (r *PaymentRepo) UpdateStatus(tx *gorm.DB, id int64, fromStatus, toStatus int, failReason string) (int64, error) {
	db := r.db
	if tx != nil {
		db = tx
	}
	updates := map[string]interface{}{
		"status": toStatus,
	}
	if failReason != "" {
		updates["fail_reason"] = failReason
	}
	if toStatus == consts.PaymentStatusSuccess {
		updates["paid_at"] = gorm.Expr("NOW()")
	}
	result := db.Model(&model.PaymentRecord{}).
		Where("id = ? AND status = ?", id, fromStatus).
		Updates(updates)
	return result.RowsAffected, result.Error
}
