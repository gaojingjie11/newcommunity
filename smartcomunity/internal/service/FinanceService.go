package service

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
	"smartcommunity/pkg/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FinanceService struct{}

const (
	PayTypeOrder            = 1
	PayTypePropertyFee      = 2
	TransactionTypeTopUp    = 3
	TransactionTypeTransfer = 4
	GreenPointsPerYuan      = 10
	CentsPerGreenPoint      = 100 / GreenPointsPerYuan

	AuthTypePassword = "password"
	AuthTypeFace     = "face"
)

type MixedPaymentResult struct {
	BusinessID           int64   `json:"business_id"`
	PayType              int     `json:"pay_type"`
	TotalAmount          float64 `json:"total_amount"`
	UsedPoints           int     `json:"used_points"`
	UsedBalance          float64 `json:"used_balance"`
	RemainingGreenPoints int     `json:"remaining_green_points"`
	RemainingBalance     float64 `json:"remaining_balance"`
}

func (s *FinanceService) UnifiedPay(userID int64, businessID int64, payType int, password string) (*MixedPaymentResult, error) {
	return s.UnifiedPayWithAuth(userID, businessID, payType, password, AuthTypePassword)
}

func (s *FinanceService) UnifiedPayWithAuth(userID int64, businessID int64, payType int, password string, authType string) (*MixedPaymentResult, error) {
	var result *MixedPaymentResult
	authType = strings.ToLower(strings.TrimSpace(authType))
	if authType == "" {
		authType = AuthTypePassword
	}
	if authType != AuthTypePassword && authType != AuthTypeFace {
		return nil, errors.New("unsupported auth type")
	}

	err := global.DB.Transaction(func(tx *gorm.DB) error {
		var user model.SysUser
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
			return errors.New("user not found")
		}
		if requiresPaymentPassword(payType, authType) {
			if strings.TrimSpace(password) == "" {
				return errors.New("payment password is required")
			}
			if !utils.CheckPasswordHash(password, user.Password) {
				return errors.New("invalid payment password")
			}
		}

		switch payType {
		case PayTypeOrder:
			return s.payOrder(tx, &user, businessID, &result)
		case PayTypePropertyFee:
			return s.payPropertyFee(tx, &user, businessID, &result)
		default:
			return errors.New("unsupported pay type")
		}
	})
	if err != nil {
		log.Printf("mixed payment failed, userID=%d businessID=%d payType=%d err=%v", userID, businessID, payType, err)
		return nil, err
	}

	return result, nil
}

func (s *FinanceService) payOrder(tx *gorm.DB, user *model.SysUser, orderID int64, result **MixedPaymentResult) error {
	var order model.Order
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND user_id = ?", orderID, user.ID).
		First(&order).Error; err != nil {
		return errors.New("order not found")
	}
	if order.Status != 0 {
		return errors.New("order status does not allow payment")
	}

	paymentResult, err := s.consumeGreenPointsAndBalance(tx, user, order.TotalAmount, orderID, PayTypeOrder, "mall_consume", fmt.Sprintf("Pay order %s", order.OrderNo))
	if err != nil {
		return err
	}

	now := time.Now()
	if err := tx.Model(&model.Order{}).
		Where("id = ?", order.ID).
		Updates(map[string]interface{}{
			"status":       1,
			"used_points":  paymentResult.UsedPoints,
			"used_balance": paymentResult.UsedBalance,
			"paid_at":      &now,
		}).Error; err != nil {
		return err
	}

	*result = paymentResult
	return nil
}

func (s *FinanceService) payPropertyFee(tx *gorm.DB, user *model.SysUser, feeID int64, result **MixedPaymentResult) error {
	var fee model.PropertyFee
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND user_id = ?", feeID, user.ID).
		First(&fee).Error; err != nil {
		return errors.New("property fee not found")
	}
	if fee.Status == 1 {
		return errors.New("property fee has already been paid")
	}

	paymentResult, err := s.consumeGreenPointsAndBalance(tx, user, fee.Amount, fee.ID, PayTypePropertyFee, "property_fee", fmt.Sprintf("Pay property fee %s", fee.Month))
	if err != nil {
		return err
	}

	now := time.Now()
	if err := tx.Model(&model.PropertyFee{}).
		Where("id = ?", fee.ID).
		Updates(map[string]interface{}{
			"status":       1,
			"pay_time":     &now,
			"used_points":  paymentResult.UsedPoints,
			"used_balance": paymentResult.UsedBalance,
		}).Error; err != nil {
		return err
	}

	*result = paymentResult
	return nil
}

func (s *FinanceService) consumeGreenPointsAndBalance(tx *gorm.DB, user *model.SysUser, amount float64, relatedID int64, payType int, action, remark string) (*MixedPaymentResult, error) {
	totalCents := amountToCents(amount)
	maxPointDeduction := totalCents / CentsPerGreenPoint
	pointsToUse := minInt(user.GreenPoints, maxPointDeduction)
	balanceCentsToUse := totalCents - pointsToUse*CentsPerGreenPoint

	if amountToCents(user.Balance) < balanceCentsToUse {
		return nil, errors.New("insufficient balance after green points deduction")
	}

	updates := map[string]interface{}{}
	if pointsToUse > 0 {
		updates["green_points"] = gorm.Expr("green_points - ?", pointsToUse)
	}
	if balanceCentsToUse > 0 {
		updates["balance"] = gorm.Expr("balance - ?", centsToAmount(balanceCentsToUse))
	}
	if len(updates) > 0 {
		if err := tx.Model(&model.SysUser{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	if pointsToUse > 0 {
		record := model.GreenPointRecord{
			UserID:    user.ID,
			Action:    action,
			Points:    -pointsToUse,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&record).Error; err != nil {
			return nil, err
		}
	}

	if balanceCentsToUse > 0 {
		transaction := model.SysTransaction{
			UserID:    user.ID,
			Type:      payType,
			Amount:    -centsToAmount(balanceCentsToUse),
			RelatedID: relatedID,
			Remark:    remark,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return nil, err
		}
	}

	user.GreenPoints -= pointsToUse
	user.Balance = centsToAmount(amountToCents(user.Balance) - balanceCentsToUse)

	return &MixedPaymentResult{
		BusinessID:           relatedID,
		PayType:              payType,
		TotalAmount:          centsToAmount(totalCents),
		UsedPoints:           pointsToUse,
		UsedBalance:          centsToAmount(balanceCentsToUse),
		RemainingGreenPoints: user.GreenPoints,
		RemainingBalance:     user.Balance,
	}, nil
}

func (s *FinanceService) GetPropertyFeeList(userID int64, page, size int) ([]model.PropertyFee, int64, error) {
	var list []model.PropertyFee
	var total int64
	db := global.DB.Model(&model.PropertyFee{}).Where("user_id = ?", userID)
	db.Count(&total)

	offset := (page - 1) * size
	err := db.Order("id desc").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}

func (s *FinanceService) Recharge(userID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("recharge amount must be greater than 0")
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.SysUser{}).Where("id = ?", userID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		transaction := model.SysTransaction{
			UserID:    userID,
			Type:      TransactionTypeTopUp,
			Amount:    amount,
			Remark:    "Balance recharge",
			CreatedAt: time.Now(),
		}
		return tx.Create(&transaction).Error
	})
}

func (s *FinanceService) Transfer(fromUserID int64, targetMobile string, amount float64) error {
	if amount <= 0 {
		return errors.New("transfer amount must be greater than 0")
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		var fromUser model.SysUser
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&fromUser, fromUserID).Error; err != nil {
			return errors.New("payer not found")
		}
		if fromUser.Balance < amount {
			return errors.New("insufficient balance")
		}

		var toUser model.SysUser
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("mobile = ?", targetMobile).
			First(&toUser).Error; err != nil {
			return errors.New("payee not found")
		}
		if toUser.ID == fromUserID {
			return errors.New("cannot transfer to self")
		}

		if err := tx.Model(&model.SysUser{}).Where("id = ?", fromUserID).
			Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.SysUser{}).Where("id = ?", toUser.ID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		outgoing := model.SysTransaction{
			UserID:    fromUserID,
			Type:      TransactionTypeTransfer,
			Amount:    -amount,
			Remark:    "Transfer to " + toUser.Username,
			CreatedAt: time.Now(),
		}
		incoming := model.SysTransaction{
			UserID:    toUser.ID,
			Type:      TransactionTypeTransfer,
			Amount:    amount,
			Remark:    "Transfer from " + fromUser.Username,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&outgoing).Error; err != nil {
			return err
		}
		return tx.Create(&incoming).Error
	})
}

func (s *FinanceService) GetTransactionList(userID int64, page, size int) ([]model.SysTransaction, int64, error) {
	var list []model.SysTransaction
	var total int64
	db := global.DB.Model(&model.SysTransaction{}).Where("user_id = ?", userID)
	db.Count(&total)

	offset := (page - 1) * size
	err := db.Order("created_at desc").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}

func (s *FinanceService) CreatePropertyFee(fee *model.PropertyFee) error {
	var count int64
	global.DB.Model(&model.PropertyFee{}).
		Where("user_id = ? AND month = ?", fee.UserID, fee.Month).
		Count(&count)
	if count > 0 {
		return errors.New("property fee for this month already exists")
	}

	fee.Status = 0
	fee.UsedPoints = 0
	fee.UsedBalance = 0
	return global.DB.Create(fee).Error
}

func (s *FinanceService) ListAllPropertyFees(page, size int) ([]model.PropertyFee, int64, error) {
	var list []model.PropertyFee
	var total int64
	db := global.DB.Model(&model.PropertyFee{})
	db.Count(&total)

	offset := (page - 1) * size
	err := db.Order("id desc").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}

func amountToCents(amount float64) int {
	return int(math.Round(amount * 100))
}

func centsToAmount(cents int) float64 {
	return float64(cents) / 100
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func requiresPaymentPassword(payType int, authType string) bool {
	if payType != PayTypeOrder && payType != PayTypePropertyFee {
		return false
	}
	return authType == AuthTypePassword
}
