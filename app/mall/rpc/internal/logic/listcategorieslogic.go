package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCategoriesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCategoriesLogic {
	return &ListCategoriesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCategoriesLogic) ListCategories(in *mall.ListCategoriesReq) (*mall.CategoryListResp, error) {
	categories, err := l.svcCtx.CategorySvc.List()
	if err != nil {
		return nil, err
	}
	var list []*mall.CategoryInfo
	for _, c := range categories {
		list = append(list, toProtoCategory(&c))
	}
	return &mall.CategoryListResp{List: list, Total: int64(len(list))}, nil
}
