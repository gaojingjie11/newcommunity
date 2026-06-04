package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminMallBindStoreProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminMallBindStoreProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminMallBindStoreProductLogic {
	return &AdminMallBindStoreProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminMallBindStoreProductLogic) AdminMallBindStoreProduct(req *types.AdminBindStoreProductReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.BindStoreProduct(l.ctx, &mall.BindStoreProductReq{
		StoreId:   req.StoreId,
		ProductId: req.ProductId,
		Stock:     req.Stock,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
