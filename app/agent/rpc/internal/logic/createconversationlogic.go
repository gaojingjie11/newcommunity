package logic

import (
	"context"
	"time"

	"smartcommunity-microservices/app/agent/rpc/agent"
	"smartcommunity-microservices/app/agent/rpc/internal/model"
	"smartcommunity-microservices/app/agent/rpc/internal/svc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateConversationLogic {
	return &CreateConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateConversationLogic) CreateConversation(in *agent.CreateConversationReq) (*agent.ConversationInfo, error) {
	convID := uuid.NewString()
	title := in.Title
	if title == "" {
		title = "新对话"
	}
	now := time.Now()

	conv := &model.SysUserConversation{
		ID:        convID,
		UserID:    in.UserId,
		Title:     title,
		Summary:   "",
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := l.svcCtx.DB.Create(conv).Error
	if err != nil {
		l.Errorf("failed to create conversation in DB: %v", err)
		return nil, err
	}

	return &agent.ConversationInfo{
		Id:        conv.ID,
		Title:     conv.Title,
		Summary:   conv.Summary,
		CreatedAt: conv.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: conv.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

