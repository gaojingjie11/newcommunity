package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminMallProductListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminMallProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminMallProductListLogic {
	return &AdminMallProductListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminMallProductListLogic) AdminMallProductList(req *types.AdminProductListReq) (resp *types.ProductListResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.AdminListProducts(l.ctx, &mall.AdminListProductsReq{
		Page:        req.Page,
		Size:        req.Size,
		Name:        req.Name,
		CategoryId:  req.CategoryId,
		IsPromotion: req.IsPromotion,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.ProductInfo, 0, len(rpcResp.List))
	for _, p := range rpcResp.List {
		list = append(list, toAPIProductInfo(p))
	}
	return &types.ProductListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
