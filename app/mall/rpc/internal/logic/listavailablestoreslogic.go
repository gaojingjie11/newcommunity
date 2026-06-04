package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAvailableStoresLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAvailableStoresLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAvailableStoresLogic {
	return &ListAvailableStoresLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListAvailableStoresLogic) ListAvailableStores(in *mall.ListAvailableStoresReq) (*mall.StoreListResp, error) {
	stores, err := l.svcCtx.OrderSvc.ListAvailableStores(in.UserId, in.CartIds)
	if err != nil {
		return nil, err
	}
	var list []*mall.StoreInfo
	for _, s := range stores {
		list = append(list, toProtoStore(&s.Store))
	}
	return &mall.StoreListResp{List: list, Total: int64(len(list))}, nil
}
