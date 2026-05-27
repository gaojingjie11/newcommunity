package repository

import (
	"smartcommunity-microservices/services/mall-service/internal/consts"
	"smartcommunity-microservices/services/mall-service/internal/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) WithTx(tx *gorm.DB) *OrderRepo {
	return &OrderRepo{db: tx}
}

func (r *OrderRepo) CreateOrder(tx *gorm.DB, order *model.Order) error {
	return tx.Create(order).Error
}

func (r *OrderRepo) FindByID(id int64) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Items").Preload("Items.Product").Preload("Store").First(&order, id).Error
	return &order, err
}

func (r *OrderRepo) FindByIDForUpdate(tx *gorm.DB, id int64) (*model.Order, error) {
	var order model.Order
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Items").First(&order, id).Error
	return &order, err
}

func (r *OrderRepo) ListByUser(userID int64, page, size int, status *int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := r.db.Model(&model.Order{}).Where("user_id = ?", userID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("id desc").Preload("Items").Preload("Items.Product").Preload("Store").Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// UpdateStatus performs a conditional status update. Returns RowsAffected.
// If 0, the current status does not match fromStatus.
func (r *OrderRepo) UpdateStatus(tx *gorm.DB, id int64, fromStatus, toStatus int) (int64, error) {
	db := r.db
	if tx != nil {
		db = tx
	}
	result := db.Model(&model.Order{}).
		Where("id = ? AND status = ?", id, fromStatus).
		Update("status", toStatus)
	return result.RowsAffected, result.Error
}

// MarkAsPaid performs a conditional update: pending_payment -> paid.
func (r *OrderRepo) MarkAsPaid(tx *gorm.DB, id int64, usedBalance int64) (int64, error) {
	now := time.Now()
	result := tx.Model(&model.Order{}).
		Where("id = ? AND status = ?", id, consts.OrderStatusPendingPayment).
		Updates(map[string]interface{}{
			"status":       consts.OrderStatusPaid,
			"used_balance": usedBalance,
			"paid_at":      &now,
			"version":      gorm.Expr("version + 1"),
		})
	return result.RowsAffected, result.Error
}

// MarkAsCancelled performs a conditional update to cancelled status.
func (r *OrderRepo) MarkAsCancelled(tx *gorm.DB, id int64, fromStatus int, reason string) (int64, error) {
	now := time.Now()
	result := tx.Model(&model.Order{}).
		Where("id = ? AND status = ?", id, fromStatus).
		Updates(map[string]interface{}{
			"status":        consts.OrderStatusCancelled,
			"cancel_reason": reason,
			"cancelled_at":  &now,
			"version":       gorm.Expr("version + 1"),
		})
	return result.RowsAffected, result.Error
}

// FindExpiredPendingOrders returns pending_payment orders past their expire_at.
func (r *OrderRepo) FindExpiredPendingOrders(limit int) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("Items").
		Where("status = ? AND expire_at IS NOT NULL AND expire_at < ?", consts.OrderStatusPendingPayment, time.Now()).
		Limit(limit).Find(&orders).Error
	return orders, err
}

func (r *OrderRepo) ListAll(page, size int, status *int, keyword string) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := r.db.Model(&model.Order{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("order_no LIKE ?", like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("id desc").Preload("Items").Preload("Items.Product").Preload("Store").Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}
