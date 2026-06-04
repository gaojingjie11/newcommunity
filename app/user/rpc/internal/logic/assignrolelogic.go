package logic

import (
	"context"
	"fmt"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRoleLogic {
	return &AssignRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AssignRoleLogic) AssignRole(in *user.AssignRoleReq) (*user.BaseResp, error) {
	roles, err := l.svcCtx.AdminService.ListRoles()
	if err != nil {
		return nil, err
	}

	var roleCode string
	for _, r := range roles {
		if r.ID == in.RoleId {
			roleCode = r.Code
			break
		}
	}

	if roleCode == "" {
		return nil, fmt.Errorf("role not found for id: %d", in.RoleId)
	}

	err = l.svcCtx.AdminService.AssignRole(in.UserId, roleCode)
	if err != nil {
		return nil, err
	}

	return &user.BaseResp{
		Code:    0,
		Message: "角色分配成功",
	}, nil
}
