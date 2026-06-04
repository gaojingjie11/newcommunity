package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallCheckFavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallCheckFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallCheckFavoriteLogic {
	return &MallCheckFavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallCheckFavoriteLogic) MallCheckFavorite(req *types.CheckFavoriteReq) (resp *types.CheckFavoriteResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.CheckFavorite(l.ctx, &mall.CheckFavoriteReq{
		UserId:    getUserIDFromCtx(l.ctx),
		ProductId: req.ProductId,
	})
	if err != nil {
		return nil, err
	}
	return &types.CheckFavoriteResp{
		IsFavorite: rpcResp.IsFavorite,
	}, nil
}
