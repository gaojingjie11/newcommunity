package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRolePermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRolePermissionsLogic {
	return &GetRolePermissionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRolePermissionsLogic) GetRolePermissions(in *user.GetRolePermissionsReq) (*user.PermissionListResp, error) {
	perms, err := l.svcCtx.AdminService.GetRolePermissions(in.RoleId)
	if err != nil {
		return nil, err
	}

	var list []string
	for _, p := range perms {
		list = append(list, p.Code)
	}

	return &user.PermissionListResp{
		Permissions: list,
	}, nil
}
