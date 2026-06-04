package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchProductsLogic {
	return &SearchProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchProductsLogic) SearchProducts(in *mall.SearchProductsReq) (*mall.ProductListResp, error) {
	products, total, err := l.svcCtx.ProductSvc.Search(in.Keyword, int(in.Page), int(in.Size))
	if err != nil {
		return nil, err
	}
	var list []*mall.ProductInfo
	for _, p := range products {
		list = append(list, toProtoProduct(&p))
	}
	return &mall.ProductListResp{List: list, Total: total}, nil
}
