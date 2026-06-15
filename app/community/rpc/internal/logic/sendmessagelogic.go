package logic

import (
	"context"
	"encoding/json"
	"time"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	redis "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMessageLogic) SendMessage(in *community.SendMessageReq) (*community.BaseResp, error) {
	msg, err := l.svcCtx.MessageRepo.Create(in.UserId, in.Content)
	if err != nil {
		return nil, err
	}

	// Update Redis cache
	if l.svcCtx.Redis != nil {
		recentCacheKey := "community:messages:recent"
		totalCacheKey := "community:messages:total"

		username := ""
		avatar := ""
		if msg.User != nil {
			username = msg.User.Username
			avatar = msg.User.Avatar
		}

		info := &community.MessageInfo{
			Id:        msg.ID,
			UserId:    msg.UserID,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt.Format("2006-01-02 15:04:05"),
			Username:  username,
			Avatar:    avatar,
		}

		jsonBytes, err := json.Marshal(info)
		if err == nil {
			pipe := l.svcCtx.Redis.Pipeline()
			pipe.ZAdd(l.ctx, recentCacheKey, redis.Z{
				Score:  float64(msg.CreatedAt.Unix()),
				Member: string(jsonBytes),
			})
			pipe.Incr(l.ctx, totalCacheKey)
			pipe.Expire(l.ctx, recentCacheKey, 48*time.Hour)
			pipe.Expire(l.ctx, totalCacheKey, 48*time.Hour)
			pipe.ZRemRangeByRank(l.ctx, recentCacheKey, 0, -501)
			_, _ = pipe.Exec(l.ctx)
		}
	}

	return &community.BaseResp{Code: 0, Message: "success"}, nil
}
