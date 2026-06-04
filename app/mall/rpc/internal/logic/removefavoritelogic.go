package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveFavoriteLogic {
	return &RemoveFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveFavoriteLogic) RemoveFavorite(in *mall.RemoveFavoriteReq) (*mall.BaseResp, error) {
	err := l.svcCtx.FavoriteSvc.Remove(in.UserId, in.ProductId)
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}
