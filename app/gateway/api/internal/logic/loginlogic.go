package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		Mobile:    req.Mobile,
		Password:  req.Password,
		ClientIp:  req.ClientIp,
		UserAgent: req.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		Token: rpcResp.Token,
		UserInfo: types.UserInfo{
			Id:             rpcResp.UserInfo.Id,
			Username:       rpcResp.UserInfo.Username,
			RealName:       rpcResp.UserInfo.RealName,
			Mobile:         rpcResp.UserInfo.Mobile,
			Avatar:         rpcResp.UserInfo.Avatar,
			GreenPoints:    rpcResp.UserInfo.GreenPoints,
			Role:           rpcResp.UserInfo.Role,
			Status:         rpcResp.UserInfo.Status,
			FaceRegistered: rpcResp.UserInfo.FaceRegistered,
			FaceImageUrl:   rpcResp.UserInfo.FaceImageUrl,
			Balance:        rpcResp.UserInfo.Balance,
		},
		IsNewUser:        rpcResp.IsNewUser,
		ProfileCompleted: rpcResp.ProfileCompleted,
	}, nil
}
