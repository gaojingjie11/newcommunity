package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserRolesLogic {
	return &GetUserRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserRolesLogic) GetUserRoles(in *user.UserIDReq) (*user.GetUserRolesResp, error) {
	roles, err := l.svcCtx.AdminService.GetUserRoles(in.UserId)
	if err != nil {
		return nil, err
	}

	var ids []int64
	for _, r := range roles {
		ids = append(ids, r.ID)
	}

	return &user.GetUserRolesResp{
		RoleIds: ids,
	}, nil
}
