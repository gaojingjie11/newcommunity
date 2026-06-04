package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallFavoriteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallFavoriteListLogic {
	return &MallFavoriteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallFavoriteListLogic) MallFavoriteList() (resp *types.FavoriteListResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.ListFavorites(l.ctx, &mall.ListFavoritesReq{
		UserId: getUserIDFromCtx(l.ctx),
		Page:   1,
		Size:   1000,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.FavoriteInfo, 0, len(rpcResp.List))
	for _, f := range rpcResp.List {
		list = append(list, toAPIFavoriteInfo(f))
	}
	return &types.FavoriteListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
