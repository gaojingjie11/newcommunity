package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductLogic {
	return &UpdateProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProductLogic) UpdateProduct(in *mall.UpdateProductReq) (*mall.BaseResp, error) {
	if err := checkProductAccess(l.ctx, l.svcCtx.DB, in.Id); err != nil {
		return nil, err
	}
	err := l.svcCtx.ProductSvc.Update(&model.Product{
		ID:            in.Id,
		CategoryID:    in.CategoryId,
		Name:          in.Name,
		Description:   in.Description,
		Price:         in.Price,
		OriginalPrice: in.OriginalPrice,
		Stock:         int(in.Stock),
		ImageURL:      in.ImageUrl,
		IsPromotion:   int(in.IsPromotion),
		Status:        int(in.Status),
	})
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}
