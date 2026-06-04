package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListStoresLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListStoresLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListStoresLogic {
	return &ListStoresLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListStoresLogic) ListStores(in *mall.ListStoresReq) (*mall.StoreListResp, error) {
	stores, total, err := l.svcCtx.StoreSvc.List(int(in.Page), int(in.Size), 0)
	if err != nil {
		return nil, err
	}
	var list []*mall.StoreInfo
	for _, s := range stores {
		list = append(list, toProtoStore(&s))
	}
	return &mall.StoreListResp{List: list, Total: total}, nil
}
