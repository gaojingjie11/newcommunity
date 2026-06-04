package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePromotionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePromotionLogic {
	return &DeletePromotionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePromotionLogic) DeletePromotion(in *mall.PromotionIDReq) (*mall.BaseResp, error) {
	err := l.svcCtx.PromotionSvc.Delete(in.Id)
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}
