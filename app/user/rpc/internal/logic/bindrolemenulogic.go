package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindRoleMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindRoleMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindRoleMenuLogic {
	return &BindRoleMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BindRoleMenuLogic) BindRoleMenu(in *user.BindRoleMenuReq) (*user.BaseResp, error) {
	err := l.svcCtx.AdminService.BindRoleMenu(in.RoleId, in.MenuIds)
	if err != nil {
		return nil, err
	}

	return &user.BaseResp{
		Code:    0,
		Message: "菜单绑定成功",
	}, nil
}
