package service

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
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

	maxImageDownloadBytes = 8 << 20   // 8MB
	targetImageBytes      = 700 << 10 // 700KB

	downloadTimeout         = 8 * time.Second
	compareConnectTimeoutMs = 3000
	compareReadTimeoutMs    = 10000
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

// CompareFace downloads both images, optimizes payload size, and uses CompareFaceAdvance.
func (s *FaceService) CompareFace(urlA, urlB string) (float32, error) {
	if s == nil || s.client == nil {
		return 0, errors.New("face client is not initialized")
	}

	imageURLA := strings.TrimSpace(urlA)
	imageURLB := strings.TrimSpace(urlB)
	if imageURLA == "" || imageURLB == "" {
		return 0, errors.New("image url cannot be empty")
	}

	imageA, err := downloadAndOptimizeImage(imageURLA)
	if err != nil {
		return 0, fmt.Errorf("download image A failed: %w", err)
	}
	imageB, err := downloadAndOptimizeImage(imageURLB)
	if err != nil {
		return 0, fmt.Errorf("download image B failed: %w", err)
	}

	req := &facebody.CompareFaceAdvanceRequest{
		ImageURLAObject: bytes.NewReader(imageA),
		ImageURLBObject: bytes.NewReader(imageB),
	}
	runtime := (&util.RuntimeOptions{}).
		SetConnectTimeout(compareConnectTimeoutMs).
		SetReadTimeout(compareReadTimeoutMs)

	resp, err := s.client.CompareFaceAdvance(req, runtime)
	if err != nil {
		return 0, fmt.Errorf("compare face failed: %w", err)
	}
	if resp == nil || resp.Body == nil || resp.Body.Data == nil || resp.Body.Data.Confidence == nil {
		return 0, errors.New("invalid compare face response")
	}

	return tea.Float32Value(resp.Body.Data.Confidence), nil
}

func downloadAndOptimizeImage(rawURL string) ([]byte, error) {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid image url: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, errors.New("image url must start with http/https")
	}

	client := &http.Client{Timeout: downloadTimeout}
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

	reader := io.LimitReader(resp.Body, maxImageDownloadBytes+1)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("empty image content")
	}
	if len(data) > maxImageDownloadBytes {
		return nil, fmt.Errorf("image is too large, max %d bytes", maxImageDownloadBytes)
	}

	optimized, err := compressImageIfNeeded(data)
	if err != nil {
		// Keep original bytes if compression fails, to avoid blocking verification.
		return data, nil
	}
	return optimized, nil
}

func compressImageIfNeeded(raw []byte) ([]byte, error) {
	if len(raw) <= targetImageBytes {
		return raw, nil
	}

	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}

	qualities := []int{85, 75, 65, 55}
	best := raw

	for _, q := range qualities {
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: q}); err != nil {
			continue
		}
		b := buf.Bytes()
		if len(b) < len(best) {
			best = append([]byte(nil), b...)
		}
		if len(best) <= targetImageBytes {
			break
		}
	}

	return best, nil
}
