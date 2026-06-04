package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindRolePermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindRolePermissionsLogic {
	return &BindRolePermissionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BindRolePermissionsLogic) BindRolePermissions(in *user.BindRolePermissionsReq) (*user.BaseResp, error) {
	allPerms, err := l.svcCtx.AdminService.ListAllPermissions()
	if err != nil {
		return nil, err
	}

	codeToID := make(map[string]int64)
	for _, p := range allPerms {
		codeToID[p.Code] = p.ID
	}

	var permissionIDs []int64
	for _, code := range in.PermissionCodes {
		if id, ok := codeToID[code]; ok {
			permissionIDs = append(permissionIDs, id)
		}
	}

	err = l.svcCtx.AdminService.BindRolePermissions(in.RoleId, permissionIDs)
	if err != nil {
		return nil, err
	}

	return &user.BaseResp{
		Code:    0,
		Message: "权限绑定成功",
	}, nil
}
