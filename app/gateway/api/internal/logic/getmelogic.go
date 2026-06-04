package logic

import (
	"context"
	"errors"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMeLogic {
	return &GetMeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMeLogic) GetMe() (resp *types.UserInfo, err error) {
	userID := getUserIDFromCtx(l.ctx)
	if userID == 0 {
		return nil, errors.New("请先登录")
	}

	rpcResp, err := l.svcCtx.UserRpc.GetProfile(l.ctx, &user.UserIDReq{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	return &types.UserInfo{
		Id:             rpcResp.Id,
		Username:       rpcResp.Username,
		RealName:       rpcResp.RealName,
		Mobile:         rpcResp.Mobile,
		Avatar:         rpcResp.Avatar,
		GreenPoints:    rpcResp.GreenPoints,
		Role:           rpcResp.Role,
		Status:         rpcResp.Status,
		FaceRegistered: rpcResp.FaceRegistered,
		FaceImageUrl:   rpcResp.FaceImageUrl,
		Balance:        rpcResp.Balance,
	}, nil
}
