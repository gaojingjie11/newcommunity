package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminMallUpdateStoreProductStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminMallUpdateStoreProductStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminMallUpdateStoreProductStatusLogic {
	return &AdminMallUpdateStoreProductStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminMallUpdateStoreProductStatusLogic) AdminMallUpdateStoreProductStatus(req *types.AdminUpdateStoreProductStatusReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.UpdateStoreProductStatus(l.ctx, &mall.UpdateStoreProductStatusReq{
		StoreId:   req.StoreId,
		ProductId: req.ProductId,
		Status:    req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
