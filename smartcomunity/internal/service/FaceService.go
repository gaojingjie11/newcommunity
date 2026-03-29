package service

import (
	"errors"
	"fmt"
	"os"
	"smartcommunity/internal/config"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	facebody "github.com/alibabacloud-go/facebody-20191230/v4/client"
	"github.com/alibabacloud-go/tea/tea"
)

const (
	defaultAliyunFaceEndpoint = "facebody.cn-shanghai.aliyuncs.com"
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
	// Prefer yaml config, fallback to env vars.
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

	cfg := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyID),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}
	return facebody.NewClient(cfg)
}

// CompareFace compares two HTTPS image URLs directly on Aliyun.
func (s *FaceService) CompareFace(urlA, urlB string) (float32, error) {
	if s == nil || s.client == nil {
		return 0, errors.New("face client is not initialized")
	}

	imageURLA := strings.TrimSpace(urlA)
	imageURLB := strings.TrimSpace(urlB)
	if imageURLA == "" || imageURLB == "" {
		return 0, errors.New("image url cannot be empty")
	}

	req := &facebody.CompareFaceRequest{
		ImageURLA: tea.String(imageURLA),
		ImageURLB: tea.String(imageURLB),
	}
	resp, err := s.client.CompareFace(req)
	if err != nil {
		return 0, fmt.Errorf("compare face failed: %w", err)
	}
	if resp == nil || resp.Body == nil || resp.Body.Data == nil || resp.Body.Data.Confidence == nil {
		return 0, errors.New("invalid compare face response")
	}

	return tea.Float32Value(resp.Body.Data.Confidence), nil
}
