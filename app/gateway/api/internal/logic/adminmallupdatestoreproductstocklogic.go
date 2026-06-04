package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminMallUpdateStoreProductStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminMallUpdateStoreProductStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminMallUpdateStoreProductStockLogic {
	return &AdminMallUpdateStoreProductStockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminMallUpdateStoreProductStockLogic) AdminMallUpdateStoreProductStock(req *types.AdminUpdateStoreProductStockReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.UpdateStoreProductStock(l.ctx, &mall.UpdateStoreProductStockReq{
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
