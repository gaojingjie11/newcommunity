package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/workorder/rpc/workorderrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkorderMyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkorderMyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkorderMyListLogic {
	return &WorkorderMyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkorderMyListLogic) WorkorderMyList(req *types.ListWorkorderReq) (resp *types.WorkorderListResp, err error) {
	rpcResp, err := l.svcCtx.WorkorderRpc.ListMyWorkorders(l.ctx, &workorderrpc.ListWorkorderReq{
		UserId: getUserIDFromCtx(l.ctx),
		Status: req.Status,
		Page:   req.Page,
		Size:   req.Size,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.WorkorderInfo, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, toAPIWorkorderInfo(item))
	}
	return &types.WorkorderListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
