package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateServiceAreaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateServiceAreaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateServiceAreaLogic {
	return &CreateServiceAreaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateServiceAreaLogic) CreateServiceArea(in *mall.CreateServiceAreaReq) (*mall.BaseResp, error) {
	err := l.svcCtx.ServiceAreaSvc.Create(&model.ServiceArea{
		Name: in.Name,
	})
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}
