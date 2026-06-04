package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallAddFavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallAddFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallAddFavoriteLogic {
	return &MallAddFavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallAddFavoriteLogic) MallAddFavorite(req *types.AddFavoriteReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.AddFavorite(l.ctx, &mall.AddFavoriteReq{
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
