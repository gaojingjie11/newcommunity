package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindPromotionProductsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindPromotionProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindPromotionProductsLogic {
	return &BindPromotionProductsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BindPromotionProductsLogic) BindPromotionProducts(in *mall.BindPromotionProductsReq) (*mall.BaseResp, error) {
	err := l.svcCtx.PromotionSvc.BindProducts(in.PromotionId, in.ProductIds)
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}
