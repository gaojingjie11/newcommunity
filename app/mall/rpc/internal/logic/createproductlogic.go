package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateProductLogic) CreateProduct(in *mall.CreateProductReq) (*mall.BaseResp, error) {
	err := l.svcCtx.ProductSvc.Create(&model.Product{
		CategoryID:    in.CategoryId,
		Name:          in.Name,
		Description:   in.Description,
		Price:         in.Price,
		OriginalPrice: in.OriginalPrice,
		Stock:         int(in.Stock),
		ImageURL:      in.ImageUrl,
		IsPromotion:   int(in.IsPromotion),
		Status:        1,
	})
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}
