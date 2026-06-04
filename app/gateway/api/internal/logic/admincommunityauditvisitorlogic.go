package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCommunityAuditVisitorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCommunityAuditVisitorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCommunityAuditVisitorLogic {
	return &AdminCommunityAuditVisitorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCommunityAuditVisitorLogic) AdminCommunityAuditVisitor(req *types.AuditVisitorReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.AuditVisitor(l.ctx, &communityrpc.AuditVisitorReq{
		Id:          req.Id,
		Status:      req.Status,
		AuditRemark: req.AuditRemark,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
