package logic

import (
	"context"
	"errors"

	"smartcommunity-microservices/app/agent/rpc/agentrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveActionLogic {
	return &ApproveActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveActionLogic) ApproveAction(req *types.ApproveActionReq) (resp *types.ApproveActionResp, err error) {
	userID := getUserIDFromCtx(l.ctx)
	if userID == 0 {
		return nil, errors.New("请先登录")
	}

	rpcResp, err := l.svcCtx.AgentRpc.ApproveAction(l.ctx, &agentrpc.ApproveActionReq{
		UserId:          userID,
		ConversationId:  req.ConversationId,
		ActionId:        req.ActionId,
		PaymentPassword: req.PaymentPassword,
		FaceImageUrl:    req.FaceImageUrl,
		PayType:         req.PayType,
		ReturnUrl:       req.ReturnUrl,
	})
	if err != nil {
		return nil, err
	}

	return &types.ApproveActionResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
		Data: types.ApproveActionData{
			PayUrl:  rpcResp.PayUrl,
			OrderId: rpcResp.OrderId,
		},
	}, nil
}
