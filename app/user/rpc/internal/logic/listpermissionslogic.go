package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPermissionsLogic {
	return &ListPermissionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPermissionsLogic) ListPermissions(in *user.EmptyReq) (*user.PermissionListResp, error) {
	perms, err := l.svcCtx.AdminService.ListAllPermissions()
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
