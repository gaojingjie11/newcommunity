package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminMallCreateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminMallCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminMallCreateProductLogic {
	return &AdminMallCreateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminMallCreateProductLogic) AdminMallCreateProduct(req *types.AdminCreateProductReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.CreateProduct(l.ctx, &mall.CreateProductReq{
		CategoryId:    req.CategoryId,
		Name:          req.Name,
		Description:   req.Description,
		Price:         int64(req.Price * 100),
		OriginalPrice: int64(req.OriginalPrice * 100),
		Stock:         req.Stock,
		ImageUrl:      req.ImageUrl,
		IsPromotion:   req.IsPromotion,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
