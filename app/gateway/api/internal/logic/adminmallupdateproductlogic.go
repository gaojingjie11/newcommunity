package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminMallUpdateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminMallUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminMallUpdateProductLogic {
	return &AdminMallUpdateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminMallUpdateProductLogic) AdminMallUpdateProduct(req *types.AdminUpdateProductReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.UpdateProduct(l.ctx, &mall.UpdateProductReq{
		Id:            req.Id,
		CategoryId:    req.CategoryId,
		Name:          req.Name,
		Description:   req.Description,
		Price:         int64(req.Price * 100),
		OriginalPrice: int64(req.OriginalPrice * 100),
		Stock:         req.Stock,
		ImageUrl:      req.ImageUrl,
		IsPromotion:   req.IsPromotion,
		Status:        req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
