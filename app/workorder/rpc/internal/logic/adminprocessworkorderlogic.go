package logic

import (
	"context"
	"errors"

	"smartcommunity-microservices/app/workorder/rpc/internal/model"
	"smartcommunity-microservices/app/workorder/rpc/internal/svc"
	"smartcommunity-microservices/app/workorder/rpc/types/workorder"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminProcessWorkorderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminProcessWorkorderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminProcessWorkorderLogic {
	return &AdminProcessWorkorderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminProcessWorkorderLogic) AdminProcessWorkorder(in *workorder.ProcessWorkorderReq) (*workorder.BaseResp, error) {
	status := int(in.Status)
	if status != model.StatusProcessing && status != model.StatusCompleted && status != model.StatusRejected {
		return nil, errors.New("invalid status: must be 1 (processing), 2 (completed), or 3 (rejected)")
	}

	_, err := l.svcCtx.WorkorderRepo.Process(in.Id, in.ProcessorId, status, in.Result)
	if err != nil {
		return nil, err
	}

	return &workorder.BaseResp{Code: 0, Message: "success"}, nil
}
