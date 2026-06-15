package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallRechargeWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallRechargeWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallRechargeWalletLogic {
	return &MallRechargeWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallRechargeWalletLogic) MallRechargeWallet(req *types.RechargeWalletReq) (resp *types.RechargeWalletResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.RechargeWallet(l.ctx, &mall.RechargeWalletReq{
		UserId:    getUserIDFromCtx(l.ctx),
		Amount:    int64(req.Amount * 100),
		PayType:   req.PayType,
		ReturnUrl: req.ReturnUrl,
	})
	if err != nil {
		return nil, err
	}
	return &types.RechargeWalletResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
		PayUrl:  rpcResp.PayUrl,
	}, nil
}
