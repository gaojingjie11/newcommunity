package service

import (
	"errors"
	"fmt"

	"smartcommunity-microservices/services/mall-service/internal/consts"
	"smartcommunity-microservices/services/mall-service/internal/model"
	"smartcommunity-microservices/services/mall-service/internal/repository"

	"gorm.io/gorm"
)

type WalletService struct {
	db         *gorm.DB
	walletRepo *repository.WalletRepo
}

func NewWalletService(db *gorm.DB, walletRepo *repository.WalletRepo) *WalletService {
	return &WalletService{db: db, walletRepo: walletRepo}
}

func validateExternalDebitTransaction(tx *model.WalletTransaction, userID, amount int64, bizType, bizID string) error {
	if tx.UserID != userID || tx.Type != consts.WalletTxTypeFee || tx.Amount != -amount || tx.BizType != bizType || tx.BizID != bizID {
		return errors.New("idempotency key conflict")
	}
	return nil
}

// Recharge adds amount (in cents) to the user's wallet.
func (s *WalletService) Recharge(userID int64, amount int64, idempotencyKey string) error {
	if amount <= 0 {
		return errors.New("充值金额必须大于0")
	}
	if idempotencyKey == "" {
		return errors.New("幂等键不能为空")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		walletRepo := s.walletRepo.WithTx(tx)

		if _, err := walletRepo.GetOrCreate(userID); err != nil {
			return err
		}

		wallet, err := walletRepo.FindByUserIDForUpdate(tx, userID)
		if err != nil {
			return err
		}

		if _, err := walletRepo.Credit(tx, userID, amount); err != nil {
			return err
		}

		return walletRepo.CreateTransaction(tx, &model.WalletTransaction{
			UserID:         userID,
			Type:           consts.WalletTxTypeRecharge,
			Amount:         amount,
			BalanceBefore:  wallet.Balance,
			BalanceAfter:   wallet.Balance + amount,
			BizType:        consts.BizTypeRecharge,
			BizID:          idempotencyKey,
			IdempotencyKey: &idempotencyKey,
			Remark:         "充值",
		})
	})
}

// Transfer moves amount (in cents) from one user to another.
func (s *WalletService) Transfer(fromUserID, toUserID int64, amount int64, idempotencyKey string) error {
	if amount <= 0 {
		return errors.New("转账金额必须大于0")
	}
	if fromUserID == toUserID {
		return errors.New("不能给自己转账")
	}
	if idempotencyKey == "" {
		return errors.New("幂等键不能为空")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		walletRepo := s.walletRepo.WithTx(tx)

		if _, err := walletRepo.GetOrCreate(fromUserID); err != nil {
			return err
		}
		if _, err := walletRepo.GetOrCreate(toUserID); err != nil {
			return err
		}

		// Lock in deterministic order to avoid deadlock
		if fromUserID < toUserID {
			if _, err := walletRepo.FindByUserIDForUpdate(tx, fromUserID); err != nil {
				return err
			}
			if _, err := walletRepo.FindByUserIDForUpdate(tx, toUserID); err != nil {
				return err
			}
		} else {
			if _, err := walletRepo.FindByUserIDForUpdate(tx, toUserID); err != nil {
				return err
			}
			if _, err := walletRepo.FindByUserIDForUpdate(tx, fromUserID); err != nil {
				return err
			}
		}

		fromWallet, err := walletRepo.FindByUserIDForUpdate(tx, fromUserID)
		if err != nil {
			return err
		}

		affected, err := walletRepo.Debit(tx, fromUserID, amount)
		if err != nil {
			return err
		}
		if affected == 0 {
			return errors.New("余额不足")
		}

		toWallet, err := walletRepo.FindByUserIDForUpdate(tx, toUserID)
		if err != nil {
			return err
		}

		if _, err := walletRepo.Credit(tx, toUserID, amount); err != nil {
			return err
		}

		// Two wallet txs share one idempotency key — use suffixes to avoid UNIQUE conflict
		keyOut := idempotencyKey + ":out"
		keyIn := idempotencyKey + ":in"
		if err := walletRepo.CreateTransaction(tx, &model.WalletTransaction{
			UserID:         fromUserID,
			Type:           consts.WalletTxTypeTransfer,
			Amount:         -amount,
			BalanceBefore:  fromWallet.Balance,
			BalanceAfter:   fromWallet.Balance - amount,
			RelatedID:      toUserID,
			BizType:        consts.BizTypeTransfer,
			BizID:          idempotencyKey,
			IdempotencyKey: &keyOut,
			Remark:         fmt.Sprintf("转账给用户%d", toUserID),
		}); err != nil {
			return err
		}

		return walletRepo.CreateTransaction(tx, &model.WalletTransaction{
			UserID:         toUserID,
			Type:           consts.WalletTxTypeTransfer,
			Amount:         amount,
			BalanceBefore:  toWallet.Balance,
			BalanceAfter:   toWallet.Balance + amount,
			RelatedID:      fromUserID,
			BizType:        consts.BizTypeTransfer,
			BizID:          idempotencyKey,
			IdempotencyKey: &keyIn,
			Remark:         fmt.Sprintf("收到用户%d转账", fromUserID),
		})
	})
}

// DebitForExternal debits amount from a user's wallet for an external service (e.g. property fee).
// Returns the wallet transaction ID and balance before/after. Uses idempotency key for exactly-once semantics.
func (s *WalletService) DebitForExternal(userID, amount int64, bizType, bizID, idempotencyKey, remark string) (walletTxID, balanceBefore, balanceAfter int64, err error) {
	if amount <= 0 {
		return 0, 0, 0, errors.New("amount must be greater than 0")
	}
	if idempotencyKey == "" {
		return 0, 0, 0, errors.New("idempotency_key required")
	}
	if bizType == "" {
		bizType = consts.BizTypePropertyFee
	}

	// Idempotency check: if a transaction with this key already exists, return it
	existing, findErr := s.walletRepo.FindTransactionByIdempotencyKey(idempotencyKey)
	if findErr == nil {
		if err := validateExternalDebitTransaction(existing, userID, amount, bizType, bizID); err != nil {
			return 0, 0, 0, err
		}
		return existing.ID, existing.BalanceBefore, existing.BalanceAfter, nil
	}
	if !errors.Is(findErr, gorm.ErrRecordNotFound) {
		return 0, 0, 0, findErr
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := s.walletRepo.WithTx(tx)

		if _, err := txRepo.GetOrCreate(userID); err != nil {
			return err
		}

		wallet, err := txRepo.FindByUserIDForUpdate(tx, userID)
		if err != nil {
			return err
		}
		balanceBefore = wallet.Balance

		rows, err := txRepo.Debit(tx, userID, amount)
		if err != nil {
			return err
		}
		if rows == 0 {
			return errors.New("余额不足")
		}

		balanceAfter = balanceBefore - amount
		wtx := &model.WalletTransaction{
			UserID:         userID,
			Type:           consts.WalletTxTypeFee,
			Amount:         -amount,
			BalanceBefore:  balanceBefore,
			BalanceAfter:   balanceAfter,
			BizType:        bizType,
			BizID:          bizID,
			IdempotencyKey: &idempotencyKey,
			Remark:         remark,
		}
		if err := txRepo.CreateTransaction(tx, wtx); err != nil {
			return err
		}
		walletTxID = wtx.ID
		return nil
	})
	if err != nil {
		existing, findErr := s.walletRepo.FindTransactionByIdempotencyKey(idempotencyKey)
		if findErr == nil {
			if validateErr := validateExternalDebitTransaction(existing, userID, amount, bizType, bizID); validateErr != nil {
				return 0, 0, 0, validateErr
			}
			return existing.ID, existing.BalanceBefore, existing.BalanceAfter, nil
		}
	}
	return
}

// GetBalance returns the user's wallet balance in cents.
func (s *WalletService) GetBalance(userID int64) (int64, error) {
	wallet, err := s.walletRepo.GetOrCreate(userID)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

// ListTransactions returns paginated wallet transactions.
func (s *WalletService) ListTransactions(userID int64, page, size int, txType *int) ([]model.WalletTransaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return s.walletRepo.ListTransactions(userID, page, size, txType)
}
