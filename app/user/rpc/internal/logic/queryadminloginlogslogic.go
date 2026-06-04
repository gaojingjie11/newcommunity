package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryAdminLoginLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryAdminLoginLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryAdminLoginLogsLogic {
	return &QueryAdminLoginLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryAdminLoginLogsLogic) QueryAdminLoginLogs(in *user.QueryLoginLogsReq) (*user.QueryLoginLogsResp, error) {
	logs, total, err := l.svcCtx.LoginLogService.QueryAdminLogs(int(in.Page), int(in.Size), "", nil)
	if err != nil {
		return nil, err
	}

	var list []*user.LoginLog
	for _, logItem := range logs {
		statusStr := "失败"
		if logItem.Success {
			statusStr = "成功"
		}
		msg := logItem.FailureReason
		if msg == "" && logItem.Success {
			msg = "登录成功"
		}
		list = append(list, &user.LoginLog{
			Id:        logItem.ID,
			UserId:    logItem.AdminUserID,
			Username:  "",
			Ip:        logItem.IP,
			UserAgent: logItem.UserAgent,
			Status:    statusStr,
			Message:   msg,
			CreatedAt: logItem.LoginTime.Format("2006-01-02 15:04:05"),
		})
	}

	return &user.QueryLoginLogsResp{
		List:  list,
		Total: total,
	}, nil
}
