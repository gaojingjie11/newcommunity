package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignUserRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignUserRolesLogic {
	return &AssignUserRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AssignUserRolesLogic) AssignUserRoles(in *user.AssignUserRolesReq) (*user.BaseResp, error) {
	err := l.svcCtx.AdminService.AssignUserRoles(in.UserId, in.RoleIds)
	if err != nil {
		return nil, err
	}

	return &user.BaseResp{
		Code:    0,
		Message: "角色分配成功",
	}, nil
}
