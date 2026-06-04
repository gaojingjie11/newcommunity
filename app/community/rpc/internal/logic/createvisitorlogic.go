package logic

import (
	"context"
	"time"

	"smartcommunity-microservices/app/community/rpc/internal/model"
	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateVisitorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateVisitorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateVisitorLogic {
	return &CreateVisitorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateVisitorLogic) CreateVisitor(in *community.CreateVisitorReq) (*community.BaseResp, error) {
	releaseTime, err := time.ParseInLocation("2006-01-02 15:04:05", in.ReleaseTime, time.Local)
	if err != nil {
		releaseTime, err = time.ParseInLocation(time.RFC3339, in.ReleaseTime, time.Local)
		if err != nil {
			return nil, err
		}
	}

	validDate, err := time.ParseInLocation("2006-01-02", in.ValidDate, time.Local)
	if err != nil {
		validDate, err = time.ParseInLocation(time.RFC3339, in.ValidDate, time.Local)
		if err != nil {
			return nil, err
		}
	}

	visitor := &model.Visitor{
		UserID:       in.UserId,
		VisitorName:  in.VisitorName,
		VisitorPhone: in.VisitorPhone,
		VisitPurpose: in.VisitPurpose,
		ReleaseTime:  releaseTime,
		ValidDate:    validDate,
		Status:       0,
	}

	if err := l.svcCtx.VisitorRepo.Create(visitor); err != nil {
		return nil, err
	}

	return &community.BaseResp{Code: 0, Message: "success"}, nil
}
