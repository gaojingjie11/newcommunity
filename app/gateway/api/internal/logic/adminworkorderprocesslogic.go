package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/workorder/rpc/workorderrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminWorkorderProcessLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminWorkorderProcessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminWorkorderProcessLogic {
	return &AdminWorkorderProcessLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminWorkorderProcessLogic) AdminWorkorderProcess(req *types.ProcessWorkorderReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.WorkorderRpc.AdminProcessWorkorder(l.ctx, &workorderrpc.ProcessWorkorderReq{
		Id:          req.Id,
		ProcessorId: getUserIDFromCtx(l.ctx),
		Status:      req.Status,
		Result:      req.Result,
		Remark:      req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
