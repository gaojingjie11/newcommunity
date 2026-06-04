package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRolesLogic {
	return &ListRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListRolesLogic) ListRoles(in *user.EmptyReq) (*user.RoleListResp, error) {
	roles, err := l.svcCtx.AdminService.ListRoles()
	if err != nil {
		return nil, err
	}

	var list []*user.Role
	for _, r := range roles {
		list = append(list, &user.Role{
			Id:     r.ID,
			Name:   r.Name,
			Code:   r.Code,
			Remark: r.Remark,
		})
	}

	return &user.RoleListResp{
		Roles: list,
	}, nil
}
