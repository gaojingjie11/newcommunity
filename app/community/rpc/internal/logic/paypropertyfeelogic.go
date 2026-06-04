package logic

import (
	"context"
	"errors"
	"fmt"
	"strconv"

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
		return nil, err
	}

	return &community.BaseResp{Code: 0, Message: "success"}, nil
}
