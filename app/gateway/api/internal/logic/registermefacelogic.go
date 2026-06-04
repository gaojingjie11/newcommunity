package logic

import (
	"context"
	"errors"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterMeFaceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterMeFaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterMeFaceLogic {
	return &RegisterMeFaceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterMeFaceLogic) RegisterMeFace(req *types.RegisterFaceReq) (resp *types.BaseResp, err error) {
	userID := getUserIDFromCtx(l.ctx)
	if userID == 0 {
		return nil, errors.New("请先登录")
	}

	rpcResp, err := l.svcCtx.UserRpc.RegisterFace(l.ctx, &user.RegisterFaceReq{
		UserId:       userID,
		FaceImageUrl: req.FaceImageUrl,
	})
	if err != nil {
		return nil, err
	}

	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
