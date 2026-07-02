package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"smartcommunity-microservices/app/mall/rpc/internal/config"
	"strings"

	"github.com/smartwalle/alipay/v3"
)

type AlipayService struct {
	client    *alipay.Client
	notifyURL string
	returnURL string
}

func NewAlipayService(cfg config.Config) (*AlipayService, error) {
	if cfg.Alipay.AppID == "" || cfg.Alipay.PrivateKey == "" || cfg.Alipay.AlipayPublicKey == "" {
		return nil, errors.New("alipay credentials are not configured")
	}

	client, err := alipay.New(cfg.Alipay.AppID, cfg.Alipay.PrivateKey, !cfg.Alipay.Sandbox)
	if err != nil {
		return nil, fmt.Errorf("failed to init alipay client: %w", err)
	}

	if err = client.LoadAliPayPublicKey(cfg.Alipay.AlipayPublicKey); err != nil {
		return nil, fmt.Errorf("failed to load alipay public key: %w", err)
	}

	return &AlipayService{
		client:    client,
		notifyURL: cfg.Alipay.NotifyURL,
		returnURL: cfg.Alipay.ReturnURL,
	}, nil
}

func (s *AlipayService) GetPaymentURL(orderNo string, amountCents int64, customReturnURL string) (string, error) {
	amountYuan := fmt.Sprintf("%.2f", float64(amountCents)/100.0)

	subject := fmt.Sprintf("智慧社区商品订单 #%s", orderNo)
	if strings.HasPrefix(orderNo, "RECH_") {
		subject = fmt.Sprintf("智慧社区账户余额充值 #%s", orderNo)
	}

	var p alipay.TradePagePay
	p.NotifyURL = s.notifyURL
	if customReturnURL != "" {
		p.ReturnURL = customReturnURL
	} else {
		p.ReturnURL = s.returnURL
	}
	p.Subject = subject
	p.OutTradeNo = orderNo
	p.TotalAmount = amountYuan
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	payURL, err := s.client.TradePagePay(p)
	if err != nil {
		return "", err
	}
	return payURL.String(), nil
}

func (s *AlipayService) VerifyNotify(params map[string]string) error {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	return s.client.VerifySign(context.Background(), values)
}

func (s *AlipayService) QueryTrade(ctx context.Context, orderNo string) (*alipay.TradeQueryRsp, error) {
	if s == nil || s.client == nil {
		return nil, errors.New("alipay client not initialized")
	}
	if strings.TrimSpace(orderNo) == "" {
		return nil, errors.New("orderNo is required")
	}

	var p alipay.TradeQuery
	p.OutTradeNo = orderNo
	return s.client.TradeQuery(ctx, p)
}
