package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductsLogic {
	return &ListProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListProductsLogic) ListProducts(in *mall.ListProductsReq) (*mall.ProductListResp, error) {
	products, total, err := l.svcCtx.ProductSvc.List(int(in.Page), int(in.Size), in.CategoryId, in.Sort, in.Name)
	if err != nil {
		return nil, err
	}
	var list []*mall.ProductInfo
	for _, p := range products {
		list = append(list, toProtoProduct(&p))
	}
	return &mall.ProductListResp{List: list, Total: total}, nil
}
