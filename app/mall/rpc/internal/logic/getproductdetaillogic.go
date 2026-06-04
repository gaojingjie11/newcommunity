package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductDetailLogic {
	return &GetProductDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductDetailLogic) GetProductDetail(in *mall.ProductIDReq) (*mall.ProductDetailResp, error) {
	product, err := l.svcCtx.ProductSvc.GetDetail(in.Id)
	if err != nil {
		return nil, err
	}
	if in.UserId > 0 {
		_ = l.svcCtx.ViewLogRepo.Create(&model.ProductViewLog{
			UserID:    in.UserId,
			ProductID: in.Id,
		})
	}
	return &mall.ProductDetailResp{Product: toProtoProduct(product)}, nil
}
