package logic

import (
	"context"

	"smartcommunity-microservices/app/agent/rpc/agent"
	"smartcommunity-microservices/app/agent/rpc/internal/model"
	"smartcommunity-microservices/app/agent/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListConversationsLogic {
	return &ListConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListConversationsLogic) ListConversations(in *agent.ListConversationsReq) (*agent.ListConversationsResp, error) {
	var list []model.SysUserConversation
	err := l.svcCtx.DB.Where("user_id = ?", in.UserId).
		Order("updated_at DESC").
		Find(&list).Error
	if err != nil {
		l.Errorf("failed to list conversations for user %d: %v", in.UserId, err)
		return nil, err
	}

	var resp []*agent.ConversationInfo
	for _, conv := range list {
		resp = append(resp, &agent.ConversationInfo{
			Id:        conv.ID,
			Title:     conv.Title,
			Summary:   conv.Summary,
			CreatedAt: conv.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: conv.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &agent.ListConversationsResp{
		List: resp,
	}, nil
}

