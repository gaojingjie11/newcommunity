package logic

import (
	"context"
	"errors"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuditVisitorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuditVisitorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditVisitorLogic {
	return &AuditVisitorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AuditVisitorLogic) AuditVisitor(in *community.AuditVisitorReq) (*community.BaseResp, error) {
	if in.Status != 1 && in.Status != 2 {
		return nil, errors.New("status must be 1 (approve) or 2 (reject)")
	}

	_, err := l.svcCtx.VisitorRepo.Audit(in.Id, int(in.Status), in.AuditRemark)
	if err != nil {
		return nil, err
	}

	return &community.BaseResp{Code: 0, Message: "success"}, nil
}
