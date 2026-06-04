package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListStoreProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListStoreProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListStoreProductsLogic {
	return &ListStoreProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListStoreProductsLogic) ListStoreProducts(in *mall.ListStoreProductsReq) (*mall.StoreProductListResp, error) {
	storeProducts, err := l.svcCtx.StoreSvc.ListProducts(in.StoreId)
	if err != nil {
		return nil, err
	}
	var list []*mall.StoreProductInfo
	for _, sp := range storeProducts {
		list = append(list, toProtoStoreProduct(&sp))
	}
	return &mall.StoreProductListResp{List: list, Total: int64(len(list))}, nil
}
