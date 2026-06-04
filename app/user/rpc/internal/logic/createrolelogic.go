package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/model"
	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RBAC Roles & Permissions
func (l *CreateRoleLogic) CreateRole(in *user.CreateRoleReq) (*user.Role, error) {
	r := &model.SysRole{
		Name:   in.Name,
		Code:   in.Code,
		Remark: in.Remark,
	}
	err := l.svcCtx.AdminService.CreateRole(r)
	if err != nil {
		return nil, err
	}

	return &user.Role{
		Id:     r.ID,
		Name:   r.Name,
		Code:   r.Code,
		Remark: r.Remark,
	}, nil
}
