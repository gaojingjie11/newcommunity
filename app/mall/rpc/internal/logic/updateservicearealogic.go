package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateServiceAreaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateServiceAreaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateServiceAreaLogic {
	return &UpdateServiceAreaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateServiceAreaLogic) UpdateServiceArea(in *mall.UpdateServiceAreaReq) (*mall.BaseResp, error) {
	err := l.svcCtx.ServiceAreaSvc.Update(&model.ServiceArea{
		ID:     in.Id,
		Name:   in.Name,
		Status: int(in.Status),
	})
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}
