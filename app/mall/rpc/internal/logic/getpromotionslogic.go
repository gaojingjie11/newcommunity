package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPromotionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPromotionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPromotionsLogic {
	return &GetPromotionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPromotionsLogic) GetPromotions(in *mall.ListProductsReq) (*mall.ProductListResp, error) {
	products, total, err := l.svcCtx.ProductSvc.GetPromotions(int(in.Page), int(in.Size))
	if err != nil {
		return nil, err
	}
	var list []*mall.ProductInfo
	for _, p := range products {
		list = append(list, toProtoProduct(&p))
	}
	return &mall.ProductListResp{List: list, Total: total}, nil
}
