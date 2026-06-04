package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMembersLogic {
	return &ListMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMembersLogic) ListMembers(in *user.ListMembersReq) (*user.ListMembersResp, error) {
	users, total, err := l.svcCtx.AdminService.ListMembers(int(in.Page), int(in.Size), "")
	if err != nil {
		return nil, err
	}

	var list []*user.UserInfo
	for _, u := range users {
		list = append(list, mapUserInfo(&u))
	}

	return &user.ListMembersResp{
		List:  list,
		Total: total,
	}, nil
}
