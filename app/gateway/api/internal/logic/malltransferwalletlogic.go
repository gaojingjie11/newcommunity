package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallTransferWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallTransferWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallTransferWalletLogic {
	return &MallTransferWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallTransferWalletLogic) MallTransferWallet(req *types.TransferWalletReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.TransferWallet(l.ctx, &mall.TransferWalletReq{
		UserId:       getUserIDFromCtx(l.ctx),
		TargetMobile: req.TargetMobile,
		Amount:       int64(req.Amount * 100),
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
