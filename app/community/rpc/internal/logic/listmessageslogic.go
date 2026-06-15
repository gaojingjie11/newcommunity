package logic

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	redis "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMessagesLogic {
	return &ListMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMessagesLogic) ListMessages(in *community.ListMessagesReq) (*community.MessageListResp, error) {
	recentCacheKey := "community:messages:recent"
	totalCacheKey := "community:messages:total"

	// Only cache the first page with size <= 100
	if l.svcCtx.Redis != nil && in.Page == 1 && in.Size <= 100 {
		totalExists, err := l.svcCtx.Redis.Exists(l.ctx, totalCacheKey).Result()
		if err == nil && totalExists > 0 {
			totalStr, err := l.svcCtx.Redis.Get(l.ctx, totalCacheKey).Result()
			if err == nil {
				total, parseErr := strconv.ParseInt(totalStr, 10, 64)
				if parseErr == nil {
					start := int64(0)
					end := int64(in.Size - 1)
					vals, err := l.svcCtx.Redis.ZRevRange(l.ctx, recentCacheKey, start, end).Result()
					if err == nil {
						if len(vals) > 0 || total == 0 {
							var list []*community.MessageInfo
							for _, val := range vals {
								var msg community.MessageInfo
								if err := json.Unmarshal([]byte(val), &msg); err == nil {
									list = append(list, &msg)
								}
							}
							return &community.MessageListResp{
								List:  list,
								Total: total,
							}, nil
						}
					}
				}
			}
		}

		// Cache Miss - Query Database (Load 100 messages to fully populate page 1 cache)
		messages, total, err := l.svcCtx.MessageRepo.List(1, 100)
		if err != nil {
			return nil, err
		}

		var list []*community.MessageInfo
		for _, m := range messages {
			username := ""
			avatar := ""
			if m.User != nil {
				username = m.User.Username
				avatar = m.User.Avatar
			}
			list = append(list, &community.MessageInfo{
				Id:        m.ID,
				UserId:    m.UserID,
				Content:   m.Content,
				CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
				Username:  username,
				Avatar:    avatar,
			})
		}

		// Populate Redis cache
		pipe := l.svcCtx.Redis.Pipeline()
		pipe.Del(l.ctx, recentCacheKey)
		for _, m := range messages {
			var info *community.MessageInfo
			for _, item := range list {
				if item.Id == m.ID {
					info = item
					break
				}
			}
			if info != nil {
				jsonBytes, err := json.Marshal(info)
				if err == nil {
					pipe.ZAdd(l.ctx, recentCacheKey, redis.Z{
						Score:  float64(m.CreatedAt.Unix()),
						Member: string(jsonBytes),
					})
				}
			}
		}
		pipe.Set(l.ctx, totalCacheKey, total, 48*time.Hour)
		pipe.Expire(l.ctx, recentCacheKey, 48*time.Hour)
		pipe.ZRemRangeByRank(l.ctx, recentCacheKey, 0, -501)
		_, _ = pipe.Exec(l.ctx)

		var respList []*community.MessageInfo
		if int(in.Size) < len(list) {
			respList = list[:in.Size]
		} else {
			respList = list
		}

		return &community.MessageListResp{
			List:  respList,
			Total: total,
		}, nil
	}

	// Fallback to database query for page > 1 or size > 100 or when Redis is not available
	messages, total, err := l.svcCtx.MessageRepo.List(int(in.Page), int(in.Size))
	if err != nil {
		return nil, err
	}

	var list []*community.MessageInfo
	for _, m := range messages {
		username := ""
		avatar := ""
		if m.User != nil {
			username = m.User.Username
			avatar = m.User.Avatar
		}

		list = append(list, &community.MessageInfo{
			Id:        m.ID,
			UserId:    m.UserID,
			Content:   m.Content,
			CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
			Username:  username,
			Avatar:    avatar,
		})
	}

	return &community.MessageListResp{
		List:  list,
		Total: total,
	}, nil
}
