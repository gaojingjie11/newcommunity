package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallProductListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallProductListLogic {
	return &MallProductListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallProductListLogic) MallProductList(req *types.ListProductsReq) (resp *types.ProductListResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.ListProducts(l.ctx, &mall.ListProductsReq{
		Page:       req.Page,
		Size:       req.Size,
		CategoryId: req.CategoryId,
		Name:       req.Name,
		Sort:       req.Sort,
		MinPrice:   int64(req.MinPrice * 100),
		MaxPrice:   int64(req.MaxPrice * 100),
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
