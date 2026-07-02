package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWalletBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWalletBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWalletBalanceLogic {
	return &GetWalletBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWalletBalanceLogic) GetWalletBalance(in *mall.UserIDReq) (*mall.WalletBalanceResp, error) {
	if l.svcCtx.PaymentReconcileSvc != nil {
		l.svcCtx.PaymentReconcileSvc.ReconcilePendingRechargesForUser(l.ctx, in.UserId)
	}
	balance, err := l.svcCtx.WalletSvc.GetBalance(in.UserId)
	if err != nil {
		return nil, err
	}
	return &mall.WalletBalanceResp{Balance: balance}, nil
}
