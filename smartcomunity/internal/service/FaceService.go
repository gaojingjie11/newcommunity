package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"smartcommunity/internal/config"
	"strings"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	facebody "github.com/alibabacloud-go/facebody-20191230/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

const (
	defaultAliyunFaceEndpoint = "facebody.cn-shanghai.aliyuncs.com"
	maxImageBytes             = 20 << 20 // 20MB
)

type FaceService struct {
	client *facebody.Client
}

func NewFaceService() (*FaceService, error) {
	client, err := newAliyunFaceClient()
	if err != nil {
		return nil, err
	}
	return &FaceService{client: client}, nil
}

func newAliyunFaceClient() (*facebody.Client, error) {
	// 优先读取 yaml 配置，兼容读取环境变量作为兜底
	accessKeyID := strings.TrimSpace(config.Conf.FaceBody.AccessKeyID)
	accessKeySecret := strings.TrimSpace(config.Conf.FaceBody.AccessKeySecret)
	endpoint := strings.TrimSpace(config.Conf.FaceBody.Endpoint)

	if accessKeyID == "" {
		accessKeyID = strings.TrimSpace(os.Getenv("ALIYUN_ACCESS_KEY_ID"))
	}
	if accessKeySecret == "" {
		accessKeySecret = strings.TrimSpace(os.Getenv("ALIYUN_ACCESS_KEY_SECRET"))
	}
	if endpoint == "" {
		endpoint = strings.TrimSpace(os.Getenv("ALIYUN_FACEBODY_ENDPOINT"))
	}
	if endpoint == "" {
		endpoint = defaultAliyunFaceEndpoint
	}

	if accessKeyID == "" || accessKeySecret == "" {
		return nil, errors.New("aliyun face credentials not configured: set facebody.access_key_id / facebody.access_key_secret in config yaml")
	}

	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyID),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}
	return facebody.NewClient(config)
}

// CompareFace 使用阿里云人脸 1:1 比对接口，返回置信度分数。
func (s *FaceService) CompareFace(urlA, urlB string) (float32, error) {
	if s == nil || s.client == nil {
		return 0, errors.New("face client is not initialized")
	}

	imageURLA := strings.TrimSpace(urlA)
	imageURLB := strings.TrimSpace(urlB)
	if imageURLA == "" || imageURLB == "" {
		return 0, errors.New("image url cannot be empty")
	}

	imageA, err := downloadImageBytes(imageURLA)
	if err != nil {
		return 0, fmt.Errorf("download image A failed: %w", err)
	}
	imageB, err := downloadImageBytes(imageURLB)
	if err != nil {
		return 0, fmt.Errorf("download image B failed: %w", err)
	}

	req := &facebody.CompareFaceAdvanceRequest{
		ImageURLAObject: bytes.NewReader(imageA),
		ImageURLBObject: bytes.NewReader(imageB),
	}
	resp, err := s.client.CompareFaceAdvance(req, &util.RuntimeOptions{})
	if err != nil {
		return 0, fmt.Errorf("compare face failed: %w", err)
	}
	if resp == nil || resp.Body == nil || resp.Body.Data == nil || resp.Body.Data.Confidence == nil {
		return 0, errors.New("invalid compare face response")
	}

	return tea.Float32Value(resp.Body.Data.Confidence), nil
}

func downloadImageBytes(rawURL string) ([]byte, error) {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid image url: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, errors.New("image url must start with http/https")
	}

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "smartcommunity-face-verify/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status %d", resp.StatusCode)
	}
	if ct := strings.ToLower(resp.Header.Get("Content-Type")); ct != "" && !strings.HasPrefix(ct, "image/") {
		return nil, fmt.Errorf("content-type is not image: %s", ct)
	}

	limited := io.LimitReader(resp.Body, maxImageBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxImageBytes {
		return nil, fmt.Errorf("image is too large, max %d bytes", maxImageBytes)
	}
	if len(data) == 0 {
		return nil, errors.New("empty image content")
	}
	return data, nil
}
