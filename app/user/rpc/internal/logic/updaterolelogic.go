package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/model"
	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateRoleLogic) UpdateRole(in *user.UpdateRoleReq) (*user.BaseResp, error) {
	err := l.svcCtx.AdminService.UpdateRole(&model.SysRole{
		ID:     in.Id,
		Name:   in.Name,
		Code:   in.Code,
		Remark: in.Remark,
	})
	if err != nil {
		return nil, err
	}

	return &user.BaseResp{
		Code:    0,
		Message: "角色更新成功",
	}, nil
}
