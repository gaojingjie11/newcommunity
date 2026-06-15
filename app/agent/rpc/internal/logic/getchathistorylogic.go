package logic

import (
	"context"
	"encoding/json"

	"smartcommunity-microservices/app/agent/rpc/agent"
	"smartcommunity-microservices/app/agent/rpc/internal/model"
	"smartcommunity-microservices/app/agent/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatHistoryLogic {
	return &GetChatHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetChatHistoryLogic) GetChatHistory(in *agent.GetChatHistoryReq) (*agent.GetChatHistoryResp, error) {
	var list []model.SysUserChatMessage
	err := l.svcCtx.DB.Where("conversation_id = ? AND user_id = ?", in.ConversationId, in.UserId).
		Order("created_at ASC").
		Find(&list).Error
	if err != nil {
		l.Errorf("failed to fetch chat history for conversation %s: %v", in.ConversationId, err)
		return nil, err
	}

	var approvals []model.AgentActionApproval
	l.svcCtx.DB.Where("conversation_id = ? AND user_id = ?", in.ConversationId, in.UserId).Find(&approvals)
	type approvalState struct {
		Status        string
		ResultPayload string
	}
	approvalMap := make(map[string]approvalState)
	for _, app := range approvals {
		approvalMap[app.ID] = approvalState{
			Status:        app.Status,
			ResultPayload: app.ResultPayload,
		}
	}

	var resp []*agent.MessageInfo
	for _, msg := range list {
		var actionResolved string
		var resultPayload string
		if msg.EventType == "approval_required" && msg.EventPayload != "" {
			var payloadMap map[string]interface{}
			if errUnmarshal := json.Unmarshal([]byte(msg.EventPayload), &payloadMap); errUnmarshal == nil {
				if actionID, _ := payloadMap["action_id"].(string); actionID != "" {
					state, ok := approvalMap[actionID]
					if ok {
						status := state.Status
						if status == "executed" || status == "approved" {
							actionResolved = "approved"
							resultPayload = state.ResultPayload
						} else if status == "rejected" {
							actionResolved = "rejected"
						}
					}
				}
			}
		}

		resp = append(resp, &agent.MessageInfo{
			Id:             msg.ID,
			Role:           msg.Role,
			Content:        msg.Content,
			EventType:      msg.EventType,
			EventPayload:   msg.EventPayload,
			ActionResolved: actionResolved,
			ResultPayload:  resultPayload,
			CreatedAt:      msg.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &agent.GetChatHistoryResp{
		List: resp,
	}, nil
}

