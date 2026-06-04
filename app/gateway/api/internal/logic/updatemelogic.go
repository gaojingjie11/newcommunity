package logic

import (
	"context"
	"errors"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMeLogic {
	return &UpdateMeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMeLogic) UpdateMe(req *types.UpdateProfileReq) (resp *types.BaseResp, err error) {
	userID := getUserIDFromCtx(l.ctx)
	if userID == 0 {
		return nil, errors.New("请先登录")
	}

	rpcResp, err := l.svcCtx.UserRpc.UpdateProfile(l.ctx, &user.UpdateProfileReq{
		UserId:   userID,
		Avatar:   req.Avatar,
		Mobile:   req.Mobile,
		Username: req.Username,
		Gender:   req.Gender,
		Email:    req.Email,
		RealName: req.RealName,
		Age:      req.Age,
	})
	if err != nil {
		return nil, err
	}

	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
