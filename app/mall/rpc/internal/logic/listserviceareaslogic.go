package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListServiceAreasLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListServiceAreasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListServiceAreasLogic {
	return &ListServiceAreasLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListServiceAreasLogic) ListServiceAreas(in *mall.ListServiceAreasReq) (*mall.ServiceAreaListResp, error) {
	areas, err := l.svcCtx.ServiceAreaSvc.List()
	if err != nil {
		return nil, err
	}
	var list []*mall.ServiceAreaInfo
	for _, a := range areas {
		list = append(list, toProtoServiceArea(&a))
	}
	return &mall.ServiceAreaListResp{List: list, Total: int64(len(list))}, nil
}
