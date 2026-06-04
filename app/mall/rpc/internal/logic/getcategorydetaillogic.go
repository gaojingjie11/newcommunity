package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCategoryDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryDetailLogic {
	return &GetCategoryDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCategoryDetailLogic) GetCategoryDetail(in *mall.CategoryIDReq) (*mall.CategoryDetailResp, error) {
	category, err := l.svcCtx.CategorySvc.GetDetail(in.Id)
	if err != nil {
		return nil, err
	}
	return &mall.CategoryDetailResp{Category: toProtoCategory(category)}, nil
}
