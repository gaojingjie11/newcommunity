package model

import "time"

type SysUserConversation struct {
	ID            string    `gorm:"primaryKey;type:varchar(64)" json:"id"`
	UserID        int64     `gorm:"column:user_id;index:idx_user_conv_updated;not null" json:"user_id"`
	Title         string    `gorm:"type:varchar(128);not null;default:'新对话'" json:"title"`
	Summary       string    `gorm:"type:text" json:"summary"`
	SummaryUntil  int       `gorm:"column:summary_until;not null;default:0" json:"summary_until"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `gorm:"index:idx_user_conv_updated;not null" json:"updated_at"`
}

func (SysUserConversation) TableName() string { return "sys_user_conversation" }

type SysUserChatMessage struct {
	ID             string    `gorm:"primaryKey;type:varchar(64)" json:"id"`
	UserID         int64     `gorm:"column:user_id;index:idx_user_chat_msg;not null" json:"user_id"`
	ConversationID string    `gorm:"column:conversation_id;index:idx_user_chat_msg;type:varchar(64);not null" json:"conversation_id"`
	Role           string    `gorm:"type:varchar(32);not null" json:"role"` // user, assistant, system
	Content        string    `gorm:"type:text;not null" json:"content"`
	EventType      string    `gorm:"type:varchar(64)" json:"event_type"`
	EventPayload   string    `gorm:"type:text" json:"event_payload"`
	CreatedAt      time.Time `gorm:"index:idx_user_chat_msg_time;not null" json:"created_at"`
}

func (SysUserChatMessage) TableName() string { return "sys_user_chat_message" }

type AgentActionApproval struct {
	ID             string    `gorm:"primaryKey;type:varchar(64)" json:"id"`
	ConversationID string    `gorm:"column:conversation_id;index:idx_agent_approval;type:varchar(64);not null" json:"conversation_id"`
	UserID         int64     `gorm:"column:user_id;index:idx_agent_approval;not null" json:"user_id"`
	ActionType     string    `gorm:"type:varchar(64);not null" json:"action_type"` // create_order, pay_order, submit_repair
	RiskLevel      string    `gorm:"type:varchar(32);not null" json:"risk_level"`  // high
	ActionPayload  string    `gorm:"type:text;not null" json:"action_payload"`     // JSON args
	Status         string    `gorm:"type:varchar(32);not null" json:"status"`      // pending, approved, rejected, executed
	ResultPayload  string    `gorm:"type:text" json:"result_payload"`              // execution result json
	IdempotencyKey string    `gorm:"type:varchar(128);index" json:"idempotency_key"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (AgentActionApproval) TableName() string { return "agent_action_approval" }

