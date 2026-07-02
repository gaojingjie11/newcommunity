package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"smartcommunity-microservices/app/community/rpc/internal/repository"
	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayPropertyFeeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayPropertyFeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayPropertyFeeLogic {
	return &PayPropertyFeeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayPropertyFeeLogic) PayPropertyFee(in *community.PayPropertyFeeReq) (*community.BaseResp, error) {
	// Validate fee
	fee, err := l.svcCtx.PropertyFeeRepo.FindByID(in.Id)
	if err != nil {
		return nil, fmt.Errorf("账单不存在")
	}
	if fee.UserID != in.UserId {
		return nil, fmt.Errorf("无权支付此账单")
	}
	if fee.Status == 1 {
		return nil, repository.ErrPropertyFeePaid
	}

	// Deduct from wallet using MallRpc client (cents)
	walletKey := fmt.Sprintf("community-fee:%d", in.Id)

	debitResp, err := l.svcCtx.MallRpc.DebitWallet(l.ctx, &mall.DebitWalletReq{
		UserId:         in.UserId,
		Amount:         fee.Amount, // Stored as cents in DB
		OrderNo:        walletKey,
		BusinessType:   "property_fee",
		IdempotencyKey: walletKey,
	})
	if err != nil {
		return nil, err
	}
	if !debitResp.Success {
		return nil, errors.New("wallet debit failed")
	}

	txID, err := strconv.ParseInt(debitResp.TransactionNo, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction number: %w", err)
	}

	// Update local property fee status and insert payment record
	_, err = l.svcCtx.PropertyFeeRepo.Pay(in.UserId, in.Id, walletKey, txID)
	if err != nil {
		l.Logger.Errorf("property fee local commit failed after wallet debit, fee_id=%d, user_id=%d, wallet_tx_id=%d, err=%v", in.Id, in.UserId, txID, err)
		if _, retryErr := l.svcCtx.PropertyFeeRepo.Pay(in.UserId, in.Id, walletKey, txID); retryErr == nil || errors.Is(retryErr, repository.ErrPropertyFeePaid) {
			err = nil
		} else {
			return nil, fmt.Errorf("缴费记账失败，请稍后重试确认结果")
		}
	}

	// Publish property fee paid event to MQ
	if l.svcCtx.MQ != nil {
		event := struct {
			Event          string `json:"event"`
			FeeID          int64  `json:"fee_id"`
			UserID         int64  `json:"user_id"`
			Month          string `json:"month"`
			Amount         int64  `json:"amount"`
			PaidAt         string `json:"paid_at"`
			IdempotencyKey string `json:"idempotency_key"`
		}{
			Event:          "property_fee.paid",
			FeeID:          fee.ID,
			UserID:         fee.UserID,
			Month:          fee.Month,
			Amount:         fee.Amount,
			PaidAt:         time.Now().Format(time.RFC3339),
			IdempotencyKey: walletKey,
		}
		body, marshalErr := json.Marshal(event)
		if marshalErr == nil {
			pubCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if pubErr := l.svcCtx.MQ.PublishEvent(pubCtx, "property_fee.paid", body); pubErr != nil {
				l.Logger.Errorf("failed to publish property_fee.paid event: %v", pubErr)
			} else {
				l.Logger.Infof("published property_fee.paid event successfully for fee %d", fee.ID)
			}
		} else {
			l.Logger.Errorf("failed to marshal property_fee.paid event: %v", marshalErr)
		}
	}

	return &community.BaseResp{Code: 0, Message: "success"}, nil
}
