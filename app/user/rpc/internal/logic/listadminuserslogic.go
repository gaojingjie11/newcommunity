package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAdminUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAdminUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAdminUsersLogic {
	return &ListAdminUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Admin Operations
func (l *ListAdminUsersLogic) ListAdminUsers(in *user.ListAdminUsersReq) (*user.ListAdminUsersResp, error) {
	users, total, err := l.svcCtx.AdminService.ListAdminUsers(int(in.Page), int(in.Size), "")
	if err != nil {
		return nil, err
	}

	var list []*user.UserInfo
	for _, u := range users {
		list = append(list, mapUserInfo(&u))
	}

	return &user.ListAdminUsersResp{
		List:  list,
		Total: total,
	}, nil
}
