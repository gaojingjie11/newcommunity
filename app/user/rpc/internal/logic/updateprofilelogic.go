package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/service"
	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProfileLogic) UpdateProfile(in *user.UpdateProfileReq) (*user.BaseResp, error) {
	var genderPtr *int
	if in.Gender != nil {
		g := int(*in.Gender)
		genderPtr = &g
	}
	var agePtr *int
	if in.Age != nil {
		a := int(*in.Age)
		agePtr = &a
	}

	err := l.svcCtx.UserService.UpdateProfile(in.UserId, service.UpdateProfileRequest{
		Avatar:   in.Avatar,
		Mobile:   in.Mobile,
		Username: in.Username,
		Gender:   genderPtr,
		Email:    in.Email,
		RealName: in.RealName,
		Age:      agePtr,
	})
	if err != nil {
		return nil, err
	}

	return &user.BaseResp{
		Code:    0,
		Message: "个人信息更新成功",
	}, nil
}
