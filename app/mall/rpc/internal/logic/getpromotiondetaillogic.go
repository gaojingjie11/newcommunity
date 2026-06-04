package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPromotionDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPromotionDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPromotionDetailLogic {
	return &GetPromotionDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPromotionDetailLogic) GetPromotionDetail(in *mall.PromotionIDReq) (*mall.PromotionDetailResp, error) {
	promo, err := l.svcCtx.PromotionSvc.GetDetail(in.Id)
	if err != nil {
		return nil, err
	}
	return &mall.PromotionDetailResp{Promotion: toProtoPromotion(promo)}, nil
}
