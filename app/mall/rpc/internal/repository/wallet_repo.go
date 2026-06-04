package repository

import (
	"smartcommunity-microservices/app/mall/rpc/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepo struct {
	db *gorm.DB
}

func NewWalletRepo(db *gorm.DB) *WalletRepo {
	return &WalletRepo{db: db}
}

func (r *WalletRepo) WithTx(tx *gorm.DB) *WalletRepo {
	return &WalletRepo{db: tx}
}

func (r *WalletRepo) FindByUserID(userID int64) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	return &wallet, err
}

func (r *WalletRepo) FindByUserIDForUpdate(tx *gorm.DB, userID int64) (*model.Wallet, error) {
	var wallet model.Wallet
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).First(&wallet).Error
	return &wallet, err
}

func (r *WalletRepo) Create(wallet *model.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *WalletRepo) GetOrCreate(userID int64) (*model.Wallet, error) {
	wallet, err := r.FindByUserID(userID)
	if err == gorm.ErrRecordNotFound {
		wallet = &model.Wallet{UserID: userID, Balance: 0}
		if err := r.Create(wallet); err != nil {
			return nil, err
		}
		return wallet, nil
	}
	return wallet, err
}

// Credit atomically adds amount to balance. Uses the caller's tx.
func (r *WalletRepo) Credit(tx *gorm.DB, userID int64, amount int64) (int64, error) {
	result := tx.Model(&model.Wallet{}).Where("user_id = ?", userID).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount))
	return result.RowsAffected, result.Error
}

// Debit atomically subtracts amount from balance with overdraft guard.
// Uses the caller's tx. Returns RowsAffected; 0 means insufficient balance.
func (r *WalletRepo) Debit(tx *gorm.DB, userID int64, amount int64) (int64, error) {
	result := tx.Model(&model.Wallet{}).
		Where("user_id = ? AND balance >= ?", userID, amount).
		UpdateColumn("balance", gorm.Expr("balance - ?", amount))
	return result.RowsAffected, result.Error
}

// CreateTransaction inserts a wallet transaction record. Uses the caller's tx.
func (r *WalletRepo) CreateTransaction(tx *gorm.DB, wtx *model.WalletTransaction) error {
	return tx.Create(wtx).Error
}

// FindTransactionByIdempotencyKey looks up a wallet transaction by its idempotency key.
func (r *WalletRepo) FindTransactionByIdempotencyKey(key string) (*model.WalletTransaction, error) {
	var wtx model.WalletTransaction
	err := r.db.Where("idempotency_key = ?", key).First(&wtx).Error
	return &wtx, err
}

func (r *WalletRepo) ListTransactions(userID int64, page, size int, txType *int) ([]model.WalletTransaction, int64, error) {
	var txs []model.WalletTransaction
	var total int64

	query := r.db.Model(&model.WalletTransaction{}).Where("user_id = ?", userID)
	if txType != nil {
		query = query.Where("type = ?", *txType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("id desc").Find(&txs).Error; err != nil {
		return nil, 0, err
	}
	return txs, total, nil
}
