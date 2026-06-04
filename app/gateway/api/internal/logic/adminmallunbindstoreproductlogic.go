package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminMallUnbindStoreProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminMallUnbindStoreProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminMallUnbindStoreProductLogic {
	return &AdminMallUnbindStoreProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminMallUnbindStoreProductLogic) AdminMallUnbindStoreProduct(req *types.AdminUnbindStoreProductReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.UnbindStoreProduct(l.ctx, &mall.UnbindStoreProductReq{
		StoreId:   req.StoreId,
		ProductId: req.ProductId,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
