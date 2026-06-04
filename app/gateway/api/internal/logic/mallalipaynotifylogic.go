package logic

import (
	"context"
	"net/http"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallAlipayNotifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallAlipayNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallAlipayNotifyLogic {
	return &MallAlipayNotifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallAlipayNotifyLogic) MallAlipayNotify(w http.ResponseWriter, r *http.Request) error {
	// 1. Parse POST Form
	if err := r.ParseForm(); err != nil {
		l.Errorf("解析支付宝回调 Form 表单失败: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("fail"))
		return err
	}

	// 2. Convert to map[string]string for RPC call
	params := make(map[string]string)
	for k, v := range r.PostForm {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	// 3. Call MallRpc AlipayNotify
	resp, err := l.svcCtx.MallRpc.AlipayNotify(l.ctx, &mall.AlipayNotifyReq{
		Params: params,
	})
	if err != nil || !resp.Success {
		l.Errorf("MallRpc.AlipayNotify 失败: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("fail"))
		return err
	}

	// 4. Respond to Alipay with plain text "success"
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("success"))
	l.Infof("支付宝回调交易处理成功, out_trade_no=%s", params["out_trade_no"])
	return nil
}
