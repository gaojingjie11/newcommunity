package logic

import (
	"context"
	"time"

	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePromotionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePromotionLogic {
	return &UpdatePromotionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePromotionLogic) UpdatePromotion(in *mall.UpdatePromotionReq) (*mall.BaseResp, error) {
	startTime, _ := time.Parse("2006-01-02 15:04:05", in.StartDate)
	endTime, _ := time.Parse("2006-01-02 15:04:05", in.EndDate)
	err := l.svcCtx.PromotionSvc.Update(&model.Promotion{
		ID:        in.Id,
		Title:     in.Title,
		Type:      int(in.Type),
		StartDate: startTime,
		EndDate:   endTime,
		Status:    int(in.Status),
	})
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}
