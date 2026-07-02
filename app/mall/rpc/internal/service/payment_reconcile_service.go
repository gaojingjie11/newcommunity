package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"smartcommunity-microservices/app/mall/rpc/internal/consts"
	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/repository"

	"github.com/smartwalle/alipay/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentReconcileService struct {
	db          *gorm.DB
	paymentRepo *repository.PaymentRepo
	walletSvc   *WalletService
	alipaySvc   *AlipayService
	eventBus    *EventBus
	log         *slog.Logger
	stopCh      chan struct{}
	doneCh      chan struct{}
	startOnce   sync.Once
	stopOnce    sync.Once
}

func NewPaymentReconcileService(
	db *gorm.DB,
	paymentRepo *repository.PaymentRepo,
	walletSvc *WalletService,
	alipaySvc *AlipayService,
	eventBus *EventBus,
	log *slog.Logger,
) *PaymentReconcileService {
	return &PaymentReconcileService{
		db:          db,
		paymentRepo: paymentRepo,
		walletSvc:   walletSvc,
		alipaySvc:   alipaySvc,
		eventBus:    eventBus,
		log:         log,
		stopCh:      make(chan struct{}),
		doneCh:      make(chan struct{}),
	}
}

func (s *PaymentReconcileService) Start() {
	if s == nil || s.alipaySvc == nil {
		return
	}
	s.startOnce.Do(func() {
		s.log.Info("payment reconcile service started")
		go s.run()
	})
}

func (s *PaymentReconcileService) Stop() {
	if s == nil || s.alipaySvc == nil {
		return
	}
	s.stopOnce.Do(func() {
		close(s.stopCh)
		<-s.doneCh
		s.log.Info("payment reconcile service stopped")
	})
}

func (s *PaymentReconcileService) run() {
	defer close(s.doneCh)

	s.scanPendingRecharges(context.Background())

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.scanPendingRecharges(context.Background())
		case <-s.stopCh:
			return
		}
	}
}

func (s *PaymentReconcileService) scanPendingRecharges(ctx context.Context) {
	records, err := s.paymentRepo.ListPendingRechargeRecords(50)
	if err != nil {
		s.log.Warn("scan pending recharge records failed", "error", err)
		return
	}
	for _, record := range records {
		recordCopy := record
		if err := s.reconcileRechargeRecord(ctx, &recordCopy); err != nil {
			s.log.Warn("reconcile recharge record failed", "order_no", recordCopy.OrderNo, "error", err)
		}
	}
}

func (s *PaymentReconcileService) ReconcilePendingRechargesForUser(ctx context.Context, userID int64) {
	if s == nil || s.alipaySvc == nil || userID <= 0 {
		return
	}
	records, err := s.paymentRepo.ListPendingRechargeRecordsByUser(userID, 5)
	if err != nil {
		s.log.Warn("list user pending recharge records failed", "user_id", userID, "error", err)
		return
	}
	for _, record := range records {
		recordCopy := record
		if err := s.reconcileRechargeRecord(ctx, &recordCopy); err != nil {
			s.log.Warn("reconcile user pending recharge failed", "user_id", userID, "order_no", recordCopy.OrderNo, "error", err)
		}
	}
}

func (s *PaymentReconcileService) ConfirmRechargeSuccess(orderNo string) (bool, int64, int64, error) {
	if s == nil {
		return false, 0, 0, errors.New("payment reconcile service unavailable")
	}
	if !strings.HasPrefix(orderNo, "RECH_") {
		return false, 0, 0, fmt.Errorf("invalid recharge order no: %s", orderNo)
	}

	var (
		applied bool
		userID  int64
		amount  int64
	)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		paymentRepo := s.paymentRepo.WithTx(tx)

		var record model.PaymentRecord
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_no = ?", orderNo).First(&record).Error; err != nil {
			return err
		}

		userID = record.UserID
		amount = record.Amount

		if record.Status == consts.PaymentStatusSuccess {
			return nil
		}
		if record.Status != consts.PaymentStatusInit {
			return fmt.Errorf("payment record %s status is %d, not reconcilable", orderNo, record.Status)
		}

		affected, err := paymentRepo.UpdateStatus(tx, record.ID, consts.PaymentStatusInit, consts.PaymentStatusSuccess, "")
		if err != nil {
			return err
		}
		if affected == 0 {
			current, findErr := paymentRepo.FindByOrderNo(orderNo)
			if findErr == nil && current.Status == consts.PaymentStatusSuccess {
				return nil
			}
			return errors.New("更新充值支付记录状态失败")
		}

		if err := s.walletSvc.RechargeTx(tx, record.UserID, record.Amount, orderNo); err != nil {
			return fmt.Errorf("增加用户钱包余额失败: %w", err)
		}

		applied = true
		return nil
	})
	if err != nil {
		return false, 0, 0, err
	}

	if applied && s.eventBus != nil {
		s.eventBus.PublishWalletRecharged(userID, amount, orderNo)
	}
	return applied, userID, amount, nil
}

func (s *PaymentReconcileService) reconcileRechargeRecord(ctx context.Context, record *model.PaymentRecord) error {
	if s == nil || s.alipaySvc == nil || record == nil {
		return nil
	}
	if record.Status == consts.PaymentStatusSuccess || !strings.HasPrefix(record.OrderNo, "RECH_") {
		return nil
	}

	rsp, err := s.alipaySvc.QueryTrade(ctx, record.OrderNo)
	if err != nil {
		return err
	}
	if rsp == nil {
		return errors.New("alipay trade query returned nil response")
	}
	if rsp.IsFailure() {
		if rsp.SubCode == "ACQ.TRADE_NOT_EXIST" {
			return nil
		}
		return fmt.Errorf("alipay trade query failed: %s %s", rsp.Code, rsp.SubMsg)
	}

	switch rsp.TradeStatus {
	case alipay.TradeStatusSuccess, alipay.TradeStatusFinished:
		_, _, _, err = s.ConfirmRechargeSuccess(record.OrderNo)
		return err
	case alipay.TradeStatusClosed:
		_, err = s.paymentRepo.UpdateStatus(nil, record.ID, consts.PaymentStatusInit, consts.PaymentStatusFailed, "支付宝交易已关闭")
		return err
	default:
		return nil
	}
}
