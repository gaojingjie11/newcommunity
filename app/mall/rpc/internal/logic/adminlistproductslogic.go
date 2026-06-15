package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminListProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListProductsLogic {
	return &AdminListProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminListProductsLogic) AdminListProducts(in *mall.AdminListProductsReq) (*mall.ProductListResp, error) {
	var isPromotion *bool
	if in.IsPromotion {
		val := true
		isPromotion = &val
	}

	products, total, err := l.svcCtx.ProductSvc.AdminList(int(in.Page), int(in.Size), in.Name, in.CategoryId, isPromotion)
	if err != nil {
		return nil, err
	}

	var list []*mall.ProductInfo
	for _, p := range products {
		list = append(list, toProtoProduct(&p))
	}
	return &mall.ProductListResp{List: list, Total: total}, nil
}
