package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminMallDeleteProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminMallDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminMallDeleteProductLogic {
	return &AdminMallDeleteProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminMallDeleteProductLogic) AdminMallDeleteProduct(req *types.AdminDeleteProductReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.DeleteProduct(l.ctx, &mall.ProductIDReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
