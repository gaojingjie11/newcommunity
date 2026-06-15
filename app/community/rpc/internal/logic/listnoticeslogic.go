package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"smartcommunity-microservices/app/community/rpc/internal/model"
	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListNoticesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type NoticeListCache struct {
	Notices []model.Notice `json:"notices"`
	Total   int64          `json:"total"`
}

func NewListNoticesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNoticesLogic {
	return &ListNoticesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListNoticesLogic) ListNotices(in *community.ListNoticesReq) (*community.NoticeListResp, error) {
	if l.svcCtx.Redis == nil {
		return l.dbListNotices(in)
	}

	cacheKey := fmt.Sprintf("community:notices:list:page:%d:size:%d", in.Page, in.Size)
	cachedJSON, err := l.svcCtx.Redis.Get(l.ctx, cacheKey).Result()
	if err == nil && cachedJSON != "" {
		var cache NoticeListCache
		if err := json.Unmarshal([]byte(cachedJSON), &cache); err == nil {
			var list []*community.NoticeInfo
			for _, n := range cache.Notices {
				list = append(list, &community.NoticeInfo{
					Id:        n.ID,
					Title:     n.Title,
					Content:   n.Content,
					Publisher: n.Publisher,
					ViewCount: n.ViewCount,
					Status:    int32(n.Status),
					CreatedAt: n.CreatedAt.Format("2006-01-02 15:04:05"),
				})
			}
			return &community.NoticeListResp{
				List:  list,
				Total: cache.Total,
			}, nil
		}
	}

	// Cache miss
	notices, total, err := l.svcCtx.NoticeRepo.List(int(in.Page), int(in.Size), false)
	if err != nil {
		return nil, err
	}

	// Save to cache
	cache := NoticeListCache{
		Notices: notices,
		Total:   total,
	}
	cacheJSON, err := json.Marshal(cache)
	if err == nil {
		_ = l.svcCtx.Redis.Set(l.ctx, cacheKey, string(cacheJSON), 24*time.Hour).Err()
	}

	var list []*community.NoticeInfo
	for _, n := range notices {
		list = append(list, &community.NoticeInfo{
			Id:        n.ID,
			Title:     n.Title,
			Content:   n.Content,
			Publisher: n.Publisher,
			ViewCount: n.ViewCount,
			Status:    int32(n.Status),
			CreatedAt: n.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &community.NoticeListResp{
		List:  list,
		Total: total,
	}, nil
}

func (l *ListNoticesLogic) dbListNotices(in *community.ListNoticesReq) (*community.NoticeListResp, error) {
	notices, total, err := l.svcCtx.NoticeRepo.List(int(in.Page), int(in.Size), false)
	if err != nil {
		return nil, err
	}

	var list []*community.NoticeInfo
	for _, n := range notices {
		list = append(list, &community.NoticeInfo{
			Id:        n.ID,
			Title:     n.Title,
			Content:   n.Content,
			Publisher: n.Publisher,
			ViewCount: n.ViewCount,
			Status:    int32(n.Status),
			CreatedAt: n.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &community.NoticeListResp{
		List:  list,
		Total: total,
	}, nil
}

