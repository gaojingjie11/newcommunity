package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbindStoreProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnbindStoreProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindStoreProductLogic {
	return &UnbindStoreProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnbindStoreProductLogic) UnbindStoreProduct(in *mall.UnbindStoreProductReq) (*mall.BaseResp, error) {
	if err := checkStoreAccess(l.ctx, in.StoreId); err != nil {
		return nil, err
	}
	err := l.svcCtx.StoreSvc.UnbindProduct(in.StoreId, in.ProductId)
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}
