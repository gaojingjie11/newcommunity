package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPermissionsLogic {
	return &GetUserPermissionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserPermissionsLogic) GetUserPermissions(in *user.UserIDReq) (*user.PermissionListResp, error) {
	roles, err := l.svcCtx.AdminService.GetUserRoles(in.UserId)
	if err != nil {
		return nil, err
	}

	roleIDs := make([]int64, len(roles))
	for i, r := range roles {
		roleIDs[i] = r.ID
	}

	perms, err := l.svcCtx.AdminService.GetPermissionsByRoleIDs(roleIDs)
	if err != nil {
		return nil, err
	}

	codes := make([]string, len(perms))
	for i, p := range perms {
		codes[i] = p.Code
	}

	return &user.PermissionListResp{
		Permissions: codes,
	}, nil
}
