package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminMallStoreProductsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminMallStoreProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminMallStoreProductsLogic {
	return &AdminMallStoreProductsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminMallStoreProductsLogic) AdminMallStoreProducts(req *types.AdminGetStoreProductsReq) (resp *types.StoreProductListResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.ListStoreProducts(l.ctx, &mall.ListStoreProductsReq{
		StoreId: req.StoreId,
		Page:    1,
		Size:    1000,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.StoreProductInfo, 0, len(rpcResp.List))
	for _, sp := range rpcResp.List {
		list = append(list, toAPIStoreProductInfo(sp))
	}
	return &types.StoreProductListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
