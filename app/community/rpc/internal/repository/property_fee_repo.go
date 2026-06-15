package repository

import (
	"errors"
	"time"

	"smartcommunity-microservices/app/community/rpc/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PropertyFeeRepo struct {
	db *gorm.DB
}

func NewPropertyFeeRepo(db *gorm.DB) *PropertyFeeRepo {
	return &PropertyFeeRepo{db: db}
}

func (r *PropertyFeeRepo) Create(item *model.PropertyFee) error {
	return r.db.Create(item).Error
}

func (r *PropertyFeeRepo) FindByID(id int64) (*model.PropertyFee, error) {
	var fee model.PropertyFee
	err := r.db.First(&fee, id).Error
	return &fee, err
}

func (r *PropertyFeeRepo) ListByUser(userID int64, page, size int) ([]model.PropertyFee, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	var items []model.PropertyFee
	var total int64
	q := r.db.Model(&model.PropertyFee{}).Where("user_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("month DESC, id DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error
	return items, total, err
}

func (r *PropertyFeeRepo) ListAll(page, size int) ([]model.PropertyFee, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	var items []model.PropertyFee
	var total int64
	q := r.db.Model(&model.PropertyFee{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("month DESC, id DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error
	return items, total, err
}

func (r *PropertyFeeRepo) Pay(userID, feeID int64, idempotencyKey string, walletTxID int64) (*model.PropertyFeePayment, error) {
	var payment model.PropertyFeePayment
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Idempotency check
		if idempotencyKey != "" {
			err := tx.Where("user_id = ? AND idempotency_key = ?", userID, idempotencyKey).First(&payment).Error
			if err == nil {
				return nil
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		// Lock fee row
		var fee model.PropertyFee
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", feeID, userID).
			First(&fee).Error; err != nil {
			return err
		}
		if fee.Status == 1 {
			return ErrPropertyFeePaid
		}

		// Conditional update: only mark paid if still unpaid
		now := time.Now()
		result := tx.Model(&model.PropertyFee{}).
			Where("id = ? AND status = 0", feeID).
			Updates(map[string]interface{}{"status": 1, "paid_at": now})
		if result.RowsAffected == 0 {
			return ErrPropertyFeePaid
		}

		payment = model.PropertyFeePayment{
			PropertyFeeID:       fee.ID,
			UserID:              userID,
			Amount:              fee.Amount,
			WalletTransactionID: walletTxID,
			IdempotencyKey:      idempotencyKey,
			Status:              1,
			PaidAt:              &now,
		}
		return tx.Create(&payment).Error
	})
	return &payment, err
}

func (r *PropertyFeeRepo) ListPaymentsByUser(userID int64, page, size int) ([]model.PropertyFeePayment, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	var items []model.PropertyFeePayment
	var total int64
	q := r.db.Model(&model.PropertyFeePayment{}).Where("user_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error
	return items, total, err
}

func (r *PropertyFeeRepo) ListPayments(page, size int) ([]model.PropertyFeePayment, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	var items []model.PropertyFeePayment
	var total int64
	q := r.db.Model(&model.PropertyFeePayment{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error
	return items, total, err
}
