package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterFaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterFaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterFaceLogic {
	return &RegisterFaceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterFaceLogic) RegisterFace(in *user.RegisterFaceReq) (*user.BaseResp, error) {
	err := l.svcCtx.UserService.RegisterFace(in.UserId, in.FaceImageUrl)
	if err != nil {
		return nil, err
	}

	return &user.BaseResp{
		Code:    0,
		Message: "人脸注册成功",
	}, nil
}
