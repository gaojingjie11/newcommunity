package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenusLogic {
	return &ListMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMenusLogic) ListMenus(in *user.EmptyReq) (*user.ListMenusResp, error) {
	menus, err := l.svcCtx.AdminService.ListAllMenus()
	if err != nil {
		return nil, err
	}

	var list []*user.Menu
	for _, m := range menus {
		list = append(list, &user.Menu{
			Id:        m.ID,
			ParentId:  m.ParentID,
			Name:      m.Name,
			Path:      m.Path,
			Component: m.Component,
			Sort:      int32(m.Sort),
			Type:      int32(m.Type),
		})
	}

	return &user.ListMenusResp{
		List: list,
	}, nil
}
