package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStoreProductStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateStoreProductStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStoreProductStockLogic {
	return &UpdateStoreProductStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateStoreProductStockLogic) UpdateStoreProductStock(in *mall.UpdateStoreProductStockReq) (*mall.BaseResp, error) {
	err := l.svcCtx.StoreSvc.UpdateProductStock(in.StoreId, in.ProductId, int(in.Stock))
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}
