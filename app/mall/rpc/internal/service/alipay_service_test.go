package service

import (
	"smartcommunity-microservices/app/mall/rpc/internal/config"
	"testing"
)

func TestAlipayService_GetPaymentURL(t *testing.T) {
	// 1. Setup mock configuration
	cfg := config.Config{}
	cfg.Alipay.AppID = "9021000164637459"
	cfg.Alipay.PrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAiZvboHVnEBmDWAOYBKZh/RTlxV2rwBKXyMtMvekYTky5jB9z
E7Dtn18FePAwG6u/He9ZZsAia/qK7SPwwUc94Tht+YJjz7xuvJhORmHdDhuut8U/
3PqHbEkagPAiLPUPw6MHBs/E4HP0fRdKrkNM0mmmO0I+IJmCig86U4KX/t/bhWay
aps695jmpgKN1BggIeD3CTI2p5XftHBkwh35ITmZTy1DeMBY1dSuqaRMlr5uMxdx
cfkMT66bSOaUQowXsU/Os6HssQ1z0tpVbdJwlx3FmbEzc4gl+elwNk61BrdPFmj3
J1siCcnuXYgpEYwDNpMD4HCQFGSXXMY2638qkQIDAQABAoIBADFY8SPTtkfxvkY7
07InMJCfg96JPuQ8Rq49KaIZCxxZK1jylkQDeNNkMgQyri3eI0VK5haQ5Ecwq81q
zBWjxK8Vm2qUtdJzUorTW46l3a4Hg1pnpAVM2m+cr6J5eugAYczYk9Z/f6y2KIEL
bz6a59u1A2XQ1ZK/Oi7kUxhLhtJhUcfnnPkd6MssT5qCwfBkWp/flqZplpTGvPV1
SjuQ9wWYxBGxMFGVFof838pfK7NyCxq2Qt4gZTXwuofVz3Y/I20eI3IIrSBk34rY
1JNyTN9l5z4B5BOx7kdtY4pTaTmzV7LLylckduFim58fAJSpMRZKJwxamvj9PAix
+9Bhp4ECgYEAyowlkQmKmN1xcI3C5ls9VLUHybVILoN8ilJIEG0vaFMuj2BjhLuP
dCBi9aEH0OJz1895vnDPTeL6oDbkLP3D0LboqsvRqfPdKdMvPTuO2wiPopby0X4R
MW5KDGi9N84sAPZSG1HBM5DLA5UQe6eiX3n1QUsmfWaWID8oWDTrq+kCgYEAreyJ
0ZJEvqw8II5uLs/n5fqlpiCxtK1gRQU9fIQ3BHUS7cJgMd3J01WHmQRwlbsoEswV
dUKzllBnjXCvV+jMRH3neMeMTKlXly/2ZB2JW/s76WlA2rg/Zx8oJd7KSdxSCbe+
LgrQcDTKk+JXe1rGOCT5/BzipLBf+wFROTSnaGkCgYEAs0Qe66NaO7micU/GtEME
oTgoUGpWHHTbgUEZ7w/z6Y3Vo6hX7F5ktQ8FBwki9cm3ZcaHpfoKQJEn6S0r/nYL
HWsFukTyqEzh7eav5K4V3d5R4kFfX/MIHIvUle8NqZqcb62TNgLB0HXSeLUyBX90
wrQaUVPGGS72qEu91XPhMiECgYAikkkW2k9F43CUPBuUvIjpAviYXBlWw7vGHHOL
Y7CX9zmK/z8lymNK2c55UROb/7bIfb6qL1cJQvRCfiqse77WwnwXWvr9Zg/eIo+E
eQaLvRW8oMpeb49SzTOqy21EX0IDYn2wI0ApfaTi3nPrNjD+igMS5R78A38gorNl
fpzkOQKBgDF8PRLkqV6JR096JN9AmWu8VUJdU80lzknI1ECMfYcOoYLvzjCgQli
/DCKT6TJipbEScviA/YO7D4m8erNA1+Qgs5F1PFEtY+FeW+CiXXf/PuiDxd2UZk
dE68dpQ2sme+SF1OcT8bmr48xY00UCjkBdufzupkvJOc2bjhaRKAX4
-----END RSA PRIVATE KEY-----`
	cfg.Alipay.AlipayPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAgRiRAF3546qzFmQWdSEn
1FBkk3snWrAniye7EnYc2CVhewuTQrD2h7/33PwLs4oT/FYyBRl58ibguHvpyO1v
lkROjmLykxVwaO11K6mj2nbYKfqdEqNELkBpEZc3JV89IYVXvgY/0/PN57rV/TuM
VGGKxIzQ7a4iv5FQ26pNcsN87I1vZNrXH9I0o5wA0etehlrAJyi/FLiSDCWu0Tt9
nnJKRLE3udXnoBtogvwXosjfVA2f+vSMfrM/lRVfmf89s48M/bMKk6IkK1j2hsS8
p22dbw7S4Ucdew2syMNVPdNve7OS+nMRSAKuxHJnj+eR6rB/7RzSL6/OjVRatMD2
UQIDAQAB
-----END PUBLIC KEY-----`
	cfg.Alipay.NotifyURL = "http://101.42.34.232:8000/api/mall/payments/alipay/notify"
	cfg.Alipay.ReturnURL = "http://localhost:81/payment/result"
	cfg.Alipay.Sandbox = true

	// 2. Initialize service
	s, err := NewAlipayService(cfg)
	if err != nil {
		t.Fatalf("failed to init AlipayService: %v", err)
	}

	// 3. Generate payment URL
	payURL, err := s.GetPaymentURL("TEST_ORDER_001", 1500, "") // 15.00 Yuan
	if err != nil {
		t.Fatalf("failed to get payment URL: %v", err)
	}

	if payURL == "" {
		t.Errorf("expected non-empty payment URL")
	}

	t.Logf("Generated payment URL: %s", payURL)
}
