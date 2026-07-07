package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginCodeLogic {
	return &LoginCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginCodeLogic) LoginCode(req *types.LoginByCodeReq) (resp *types.LoginResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.LoginByCode(l.ctx, &user.LoginByCodeReq{
		Mobile:    req.Mobile,
		Code:      req.Code,
		ClientIp:  req.ClientIp,
		UserAgent: req.UserAgent,
	})
	if err != nil {
		return nil, err
	}

	// Fetch permissions
	permResp, err := l.svcCtx.UserRpc.GetUserPermissions(l.ctx, &user.UserIDReq{
		UserId: rpcResp.UserInfo.Id,
	})
	var permissions []string
	if err == nil && permResp != nil {
		permissions = permResp.Permissions
	} else {
		permissions = make([]string, 0)
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
			RoleName:       rpcResp.UserInfo.RoleName,
			Status:         rpcResp.UserInfo.Status,
			FaceRegistered: rpcResp.UserInfo.FaceRegistered,
			FaceImageUrl:   rpcResp.UserInfo.FaceImageUrl,
			Balance:        rpcResp.UserInfo.Balance,
			Gender:         rpcResp.UserInfo.Gender,
			Email:          rpcResp.UserInfo.Email,
			Age:            rpcResp.UserInfo.Age,
			Permissions:    permissions,
		},
		IsNewUser:        rpcResp.IsNewUser,
		ProfileCompleted: rpcResp.ProfileCompleted,
	}, nil
}
