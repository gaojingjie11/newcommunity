package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteServiceAreaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteServiceAreaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteServiceAreaLogic {
	return &DeleteServiceAreaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteServiceAreaLogic) DeleteServiceArea(in *mall.CategoryIDReq) (*mall.BaseResp, error) {
	err := l.svcCtx.ServiceAreaSvc.Delete(in.Id)
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}
