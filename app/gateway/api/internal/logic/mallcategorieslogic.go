package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallCategoriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallCategoriesLogic {
	return &MallCategoriesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallCategoriesLogic) MallCategories() (resp *types.CategoryListResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.ListCategories(l.ctx, &mall.ListCategoriesReq{
		Page: 1,
		Size: 1000,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.CategoryInfo, 0, len(rpcResp.List))
	for _, c := range rpcResp.List {
		list = append(list, toAPICategoryInfo(c))
	}
	return &types.CategoryListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
