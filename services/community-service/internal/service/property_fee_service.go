package service

import (
	"errors"
	"fmt"

	"smartcommunity-microservices/services/community-service/internal/client"
	"smartcommunity-microservices/services/community-service/internal/model"
	"smartcommunity-microservices/services/community-service/internal/repository"
)

type PropertyFeeService struct {
	repo       *repository.PropertyFeeRepo
	mallClient *client.MallClient
}

type CreatePropertyFeeRequest struct {
	UserID  int64  `json:"user_id" binding:"required"`
	Month   string `json:"month" binding:"required"`
	Amount  int64  `json:"amount" binding:"required"`
	DueDate string `json:"due_date"`
}

type PayPropertyFeeRequest struct {
	IdempotencyKey string `json:"idempotency_key"`
	PayType        string `json:"pay_type"`
	Password       string `json:"password"`
	FaceImageURL   string `json:"face_image_url"`
}

type PropertyFeePayResult struct {
	Payment     *model.PropertyFeePayment `json:"payment"`
	UsedPoints  int                       `json:"used_points"`
	UsedBalance float64                   `json:"used_balance"`
}

func propertyFeeWalletKey(feeID int64) string {
	return fmt.Sprintf("community-fee:%d", feeID)
}

func NewPropertyFeeService(repo *repository.PropertyFeeRepo, mallClient *client.MallClient) *PropertyFeeService {
	return &PropertyFeeService{repo: repo, mallClient: mallClient}
}

func (s *PropertyFeeService) Create(req CreatePropertyFeeRequest) (*model.PropertyFee, error) {
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	item := &model.PropertyFee{
		UserID: req.UserID,
		Month:  req.Month,
		Amount: req.Amount,
		Status: 0,
	}
	if req.DueDate != "" {
		due, err := parseDate(req.DueDate)
		if err != nil {
			return nil, err
		}
		item.DueDate = &due
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *PropertyFeeService) ListByUser(userID int64, page, size int) ([]model.PropertyFee, int64, error) {
	return s.repo.ListByUser(userID, page, size)
}

func (s *PropertyFeeService) ListAll(page, size int) ([]model.PropertyFee, int64, error) {
	return s.repo.ListAll(page, size)
}

func (s *PropertyFeeService) Pay(userID, feeID int64, req PayPropertyFeeRequest) (*PropertyFeePayResult, error) {
	walletKey := req.IdempotencyKey
	if walletKey == "" {
		walletKey = propertyFeeWalletKey(feeID)
	}
	if s.mallClient == nil {
		return nil, errors.New("mall-service wallet client is not configured")
	}

	// Step 1: Quick validation (no lock)
	fee, err := s.repo.FindByID(feeID)
	if err != nil {
		return nil, fmt.Errorf("账单不存在")
	}
	if fee.UserID != userID {
		return nil, fmt.Errorf("无权支付此账单")
	}
	if fee.Status == 1 {
		return nil, repository.ErrPropertyFeePaid
	}

	// Step 2: Call mall-service to debit wallet (idempotent)
	remark := fmt.Sprintf("物业费缴纳 %s", fee.Month)
	walletTxID, err := s.mallClient.DebitWallet(userID, fee.Amount*100, walletKey, remark, req.PayType, req.Password, req.FaceImageURL)
	if err != nil {
		return nil, err
	}
	if walletTxID <= 0 {
		return nil, errors.New("wallet transaction id is empty")
	}

	// Step 3: Lock fee → conditional update → write payment record
	payment, err := s.repo.Pay(userID, feeID, walletKey, walletTxID)
	if err != nil {
		return nil, err
	}
	return &PropertyFeePayResult{
		Payment:     payment,
		UsedPoints:  0,
		UsedBalance: float64(fee.Amount),
	}, nil
}

func (s *PropertyFeeService) ListPaymentsByUser(userID int64, page, size int) ([]model.PropertyFeePayment, int64, error) {
	return s.repo.ListPaymentsByUser(userID, page, size)
}

func (s *PropertyFeeService) ListPayments(page, size int) ([]model.PropertyFeePayment, int64, error) {
	return s.repo.ListPayments(page, size)
}
