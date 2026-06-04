package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/workorder/rpc/workorderrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkorderLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkorderLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkorderLogsLogic {
	return &WorkorderLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkorderLogsLogic) WorkorderLogs(req *types.GetLogsReq) (resp *types.LogListResp, err error) {
	rpcResp, err := l.svcCtx.WorkorderRpc.GetWorkorderLogs(l.ctx, &workorderrpc.GetLogsReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.WorkorderLogInfo, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, toAPIWorkorderLogInfo(item))
	}
	return &types.LogListResp{
		List: list,
	}, nil
}
