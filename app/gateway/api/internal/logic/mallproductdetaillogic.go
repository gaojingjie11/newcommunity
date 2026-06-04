package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallProductDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallProductDetailLogic {
	return &MallProductDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallProductDetailLogic) MallProductDetail(req *types.ProductIDReq) (resp *types.ProductInfo, err error) {
	rpcResp, err := l.svcCtx.MallRpc.GetProductDetail(l.ctx, &mall.ProductIDReq{
		Id:     req.Id,
		UserId: getUserIDFromCtx(l.ctx),
	})
	if err != nil {
		return nil, err
	}
	info := toAPIProductInfo(rpcResp.Product)
	return &info, nil
}
