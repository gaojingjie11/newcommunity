package logic

import (
	"context"
	"time"

	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePromotionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePromotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePromotionLogic {
	return &CreatePromotionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePromotionLogic) CreatePromotion(in *mall.CreatePromotionReq) (*mall.BaseResp, error) {
	startTime, _ := time.Parse("2006-01-02 15:04:05", in.StartDate)
	endTime, _ := time.Parse("2006-01-02 15:04:05", in.EndDate)
	err := l.svcCtx.PromotionSvc.Create(&model.Promotion{
		Title:     in.Title,
		Type:      int(in.Type),
		StartDate: startTime,
		EndDate:   endTime,
		Status:    1,
	})
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}
