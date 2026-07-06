package svc

import (
	"context"
	"log"

	"smartcommunity-microservices/app/agent/rpc/internal/config"
	"smartcommunity-microservices/app/agent/rpc/internal/model"
	agentservice "smartcommunity-microservices/app/agent/rpc/internal/service"
	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/mall/rpc/mallrpc"
	"smartcommunity-microservices/app/user/rpc/userrpc"
	"smartcommunity-microservices/app/workorder/rpc/workorderrpc"
	"smartcommunity-microservices/common/db"

	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config       config.Config
	DB           *gorm.DB
	UserRpc      userrpc.UserRpc
	MallRpc      mallrpc.MallRpc
	CommunityRpc communityrpc.CommunityRpc
	WorkorderRpc workorderrpc.WorkorderRpc
	KnowledgeSvc *agentservice.KnowledgeService
}

func NewServiceContext(c config.Config) *ServiceContext {
	// Initialize Postgres Connection
	database, err := db.InitPostgres(c.Postgres)
	if err != nil {
		log.Fatalf("failed to init postgres in agent-rpc: %v", err)
	}

	// Run AutoMigrate for chat history tables
	err = database.AutoMigrate(
		&model.SysUserConversation{},
		&model.SysUserChatMessage{},
		&model.AgentActionApproval{},
	)
	if err != nil {
		log.Fatalf("failed to automigrate agent-rpc GORM tables: %v", err)
	}

	var knowledgeSvc *agentservice.KnowledgeService
	if svc, err := agentservice.NewKnowledgeService(database, c.Agent); err != nil {
		log.Printf("RAG knowledge service disabled: %v", err)
	} else if err := svc.Init(context.Background()); err != nil {
		log.Printf("RAG knowledge service init failed: %v", err)
	} else {
		knowledgeSvc = svc
		knowledgeSvc.StartBackgroundSync()
	}

	return &ServiceContext{
		Config:       c,
		DB:           database,
		UserRpc:      userrpc.NewUserRpc(zrpc.MustNewClient(c.UserRpc)),
		MallRpc:      mallrpc.NewMallRpc(zrpc.MustNewClient(c.MallRpc)),
		CommunityRpc: communityrpc.NewCommunityRpc(zrpc.MustNewClient(c.CommunityRpc)),
		WorkorderRpc: workorderrpc.NewWorkorderRpc(zrpc.MustNewClient(c.WorkorderRpc)),
		KnowledgeSvc: knowledgeSvc,
	}
}
