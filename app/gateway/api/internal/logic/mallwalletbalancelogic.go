package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallWalletBalanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallWalletBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallWalletBalanceLogic {
	return &MallWalletBalanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallWalletBalanceLogic) MallWalletBalance() (resp *types.WalletBalanceResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.GetWalletBalance(l.ctx, &mall.UserIDReq{
		UserId: getUserIDFromCtx(l.ctx),
	})
	if err != nil {
		return nil, err
	}
	return &types.WalletBalanceResp{
		Balance: float64(rpcResp.Balance) / 100.0,
	}, nil
}
