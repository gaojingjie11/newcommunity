package logic

import (
	"context"
	"errors"
	"strings"

	"smartcommunity-microservices/app/workorder/rpc/internal/model"
	"smartcommunity-microservices/app/workorder/rpc/internal/svc"
	"smartcommunity-microservices/app/workorder/rpc/types/workorder"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWorkorderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateWorkorderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkorderLogic {
	return &CreateWorkorderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateWorkorderLogic) CreateWorkorder(in *workorder.CreateWorkorderReq) (*workorder.BaseResp, error) {
	wType := strings.ToLower(strings.TrimSpace(in.Type))
	if wType != model.WorkorderTypeRepair && wType != model.WorkorderTypeComplaint {
		// Try mapping user friendly terms
		if wType == "1" || wType == "报修" || wType == "设备报修" {
			wType = model.WorkorderTypeRepair
		} else if wType == "2" || wType == "投诉" || wType == "建议投诉" {
			wType = model.WorkorderTypeComplaint
		} else {
			return nil, errors.New("type must be repair or complaint")
		}
	}

	category := strings.TrimSpace(in.Category)
	description := strings.TrimSpace(in.Description)
	if category == "" || description == "" {
		return nil, errors.New("category and description required")
	}

	item := &model.WorkOrder{
		Type:        wType,
		UserID:      in.UserId,
		Category:    category,
		Description: description,
		Status:      model.StatusPending,
	}

	if err := l.svcCtx.WorkorderRepo.Create(item); err != nil {
		return nil, err
	}

	return &workorder.BaseResp{Code: 0, Message: "success"}, nil
}
