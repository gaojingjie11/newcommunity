package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFavoritesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFavoritesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFavoritesLogic {
	return &ListFavoritesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListFavoritesLogic) ListFavorites(in *mall.ListFavoritesReq) (*mall.FavoriteListResp, error) {
	favorites, total, err := l.svcCtx.FavoriteSvc.List(in.UserId, int(in.Page), int(in.Size))
	if err != nil {
		return nil, err
	}
	var list []*mall.FavoriteInfo
	for _, f := range favorites {
		list = append(list, toProtoFavorite(&f))
	}
	return &mall.FavoriteListResp{List: list, Total: total}, nil
}
