package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallDeleteFavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallDeleteFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallDeleteFavoriteLogic {
	return &MallDeleteFavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallDeleteFavoriteLogic) MallDeleteFavorite(req *types.DeleteFavoriteReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.RemoveFavorite(l.ctx, &mall.RemoveFavoriteReq{
		UserId:    getUserIDFromCtx(l.ctx),
		ProductId: req.ProductId,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
