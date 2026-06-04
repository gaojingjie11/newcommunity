package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckFavoriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckFavoriteLogic {
	return &CheckFavoriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckFavoriteLogic) CheckFavorite(in *mall.CheckFavoriteReq) (*mall.CheckFavoriteResp, error) {
	isFav, err := l.svcCtx.FavoriteSvc.Check(in.UserId, in.ProductId)
	if err != nil {
		return nil, err
	}
	return &mall.CheckFavoriteResp{IsFavorite: isFav}, nil
}
