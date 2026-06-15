package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"smartcommunity-microservices/app/mall/rpc/internal/config"

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
