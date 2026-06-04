package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommunityCreateVisitorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommunityCreateVisitorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommunityCreateVisitorLogic {
	return &CommunityCreateVisitorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommunityCreateVisitorLogic) CommunityCreateVisitor(req *types.CreateVisitorReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.CreateVisitor(l.ctx, &communityrpc.CreateVisitorReq{
		UserId:       getUserIDFromCtx(l.ctx),
		VisitorName:  req.VisitorName,
		VisitorPhone: req.VisitorPhone,
		VisitPurpose: req.VisitPurpose,
		ReleaseTime:  req.ReleaseTime,
		ValidDate:    req.ValidDate,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
