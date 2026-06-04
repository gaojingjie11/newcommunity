package logic

import (
	"context"

	"smartcommunity-microservices/app/workorder/rpc/internal/svc"
	"smartcommunity-microservices/app/workorder/rpc/types/workorder"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkorderLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWorkorderLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkorderLogsLogic {
	return &GetWorkorderLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWorkorderLogsLogic) GetWorkorderLogs(in *workorder.GetLogsReq) (*workorder.LogListResp, error) {
	logs, err := l.svcCtx.WorkorderRepo.ListLogs(in.Id)
	if err != nil {
		return nil, err
	}

	var list []*workorder.WorkorderLogInfo
	for _, log := range logs {
		list = append(list, &workorder.WorkorderLogInfo{
			Id:         log.ID,
			TargetType: log.TargetType,
			TargetId:   log.TargetID,
			FromStatus: int32(log.FromStatus),
			ToStatus:   int32(log.ToStatus),
			OperatorId: log.OperatorID,
			Action:     log.Action,
			Remark:     log.Remark,
			CreatedAt:   log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &workorder.LogListResp{
		List: list,
	}, nil
}
