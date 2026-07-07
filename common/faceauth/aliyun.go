package faceauth

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	facebody "github.com/alibabacloud-go/facebody-20191230/v6/client"
	"github.com/alibabacloud-go/tea/dara"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	defaultEndpoint         = "facebody.cn-shanghai.aliyuncs.com"
	defaultRegion           = "cn-shanghai"
	defaultMinConfidence    = float32(75)
	defaultQualityThreshold = float32(70)

	maxImageDownloadBytes = 8 << 20
	targetImageBytes      = 700 << 10
	downloadTimeout       = 8 * time.Second
	compareConnectTimeout = 3000
	compareReadTimeout    = 10000
)

type compareConfig struct {
	endpoint         string
	region           string
	minConfidence    float32
	qualityThreshold float32
	mode             string
}

type CompareResult struct {
	Confidence    float64 `json:"confidence"`
	MinConfidence float64 `json:"min_confidence"`
	QualityA      float64 `json:"quality_a"`
	QualityB      float64 `json:"quality_b"`
	MessageTips   string  `json:"message_tips"`
}

type minioResolver struct {
	client    *minio.Client
	bucket    string
	publicURL string
}

var (
	loadOnce     sync.Once
	sharedClient *facebody.Client
	sharedCfg    compareConfig
	sharedErr    error

	minioOnce      sync.Once
	sharedMinio    *minioResolver
	sharedMinioErr error
)

func ValidateEnrollment(ctx context.Context, imageURL string) error {
	_ = ctx
	imageURL = strings.TrimSpace(imageURL)
	if imageURL == "" {
		return errors.New("请先上传人脸照片")
	}
	return nil
}

func VerifyMatch(ctx context.Context, registeredURL, capturedURL string) (*CompareResult, error) {
	_ = ctx
	registeredURL = strings.TrimSpace(registeredURL)
	capturedURL = strings.TrimSpace(capturedURL)
	if registeredURL == "" {
		return nil, errors.New("当前账号未录入人脸")
	}
	if capturedURL == "" {
		return nil, errors.New("请先完成刷脸验证")
	}

	client, cfg, err := getClient()
	if err != nil {
		return nil, err
	}

	resp, err := compareFace(client, cfg, registeredURL, capturedURL)
	if err != nil {
		log.Printf("[FaceAuth] compare failed: registered=%s captured=%s err=%v", registeredURL, capturedURL, err)
		return nil, classifyCompareError(err)
	}
	if resp == nil || resp.Body == nil || resp.Body.Data == nil {
		return nil, errors.New("人脸比对结果异常，请重新尝试")
	}

	data := resp.Body.Data
	result := &CompareResult{
		MinConfidence: float64(cfg.minConfidence),
	}
	if data.Confidence != nil {
		result.Confidence = float64(*data.Confidence)
	}
	if data.QualityScoreA != nil {
		result.QualityA = float64(*data.QualityScoreA)
	}
	if data.QualityScoreB != nil {
		result.QualityB = float64(*data.QualityScoreB)
	}
	if data.MessageTips != nil {
		result.MessageTips = strings.TrimSpace(*data.MessageTips)
	}

	if result.Confidence <= 0 {
		return result, errors.New("未能识别有效人脸，请正对摄像头重新尝试")
	}
	if float32(result.Confidence) < cfg.minConfidence {
		return result, fmt.Errorf("人脸比对未通过，相似度 %.2f 低于阈值 %.2f", result.Confidence, result.MinConfidence)
	}

	return result, nil
}

func getClient() (*facebody.Client, compareConfig, error) {
	loadOnce.Do(func() {
		cfg := compareConfig{
			endpoint:         getEnvOrDefault("ALIYUN_FACEBODY_ENDPOINT", defaultEndpoint),
			region:           getEnvOrDefault("ALIYUN_FACEBODY_REGION", defaultRegion),
			minConfidence:    getEnvAsFloat32("ALIYUN_FACE_COMPARE_MIN_CONFIDENCE", defaultMinConfidence),
			qualityThreshold: getEnvAsFloat32("ALIYUN_FACE_COMPARE_QUALITY_THRESHOLD", defaultQualityThreshold),
			mode:             getCompareMode(),
		}

		accessKeyID := strings.TrimSpace(os.Getenv("ALIYUN_ACCESS_KEY_ID"))
		accessKeySecret := strings.TrimSpace(os.Getenv("ALIYUN_ACCESS_KEY_SECRET"))
		if accessKeyID == "" || accessKeySecret == "" {
			sharedErr = errors.New("人脸支付服务未配置阿里云密钥")
			return
		}

		client, err := facebody.NewClient((&openapi.Config{}).
			SetAccessKeyId(accessKeyID).
			SetAccessKeySecret(accessKeySecret).
			SetRegionId(cfg.region).
			SetEndpoint(cfg.endpoint))
		if err != nil {
			sharedErr = errors.New("初始化人脸支付服务失败")
			return
		}

		sharedClient = client
		sharedCfg = cfg
	})

	return sharedClient, sharedCfg, sharedErr
}

func getEnvOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getEnvAsFloat32(key string, fallback float32) float32 {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return fallback
	}
	return float32(parsed)
}

func formatMessageTips(raw string) string {
	lower := strings.ToLower(strings.TrimSpace(raw))
	switch {
	case strings.Contains(lower, "quality score less threshold"):
		return "当前照片清晰度或光照不足，请正对摄像头重新尝试"
	case strings.Contains(lower, "no face"):
		return "未识别到人脸，请正对摄像头重新尝试"
	case strings.Contains(lower, "mask"):
		return "请摘下口罩后重新尝试"
	default:
		return "人脸校验未通过，请确保仅本人正对摄像头并保持光线充足"
	}
}

func classifyCompareError(err error) error {
	if err == nil {
		return nil
	}
	msg := strings.ToLower(strings.TrimSpace(err.Error()))
	switch {
	case strings.Contains(msg, "registered image download failed"):
		return errors.New("登记人脸底图读取失败，请重新录入人脸后再试")
	case strings.Contains(msg, "captured image download failed"):
		return errors.New("本次抓拍图片读取失败，请重新抓拍后再试")
	case strings.Contains(msg, "http status 403"), strings.Contains(msg, "http status 404"):
		return errors.New("人脸图片访问失败，请重新抓拍或重新录入后再试")
	case strings.Contains(msg, "content-type is not image"), strings.Contains(msg, "empty image content"):
		return errors.New("人脸图片内容异常，请重新抓拍后再试")
	case strings.Contains(msg, "image is too large"):
		return errors.New("人脸图片过大，请重新抓拍后再试")
	case strings.Contains(msg, "timeout"), strings.Contains(msg, "i/o timeout"), strings.Contains(msg, "connection refused"):
		return errors.New("人脸服务连接超时，请稍后重试")
	default:
		return errors.New("人脸比对失败，请稍后重试")
	}
}

func compareFace(client *facebody.Client, cfg compareConfig, registeredURL, capturedURL string) (*facebody.CompareFaceResponse, error) {
	return compareFaceByDownload(client, cfg, registeredURL, capturedURL)
}

func compareFaceByDownload(client *facebody.Client, cfg compareConfig, registeredURL, capturedURL string) (*facebody.CompareFaceResponse, error) {
	imageA, err := downloadAndOptimizeImage(registeredURL)
	if err != nil {
		return nil, fmt.Errorf("registered image download failed: %w", err)
	}
	imageB, err := downloadAndOptimizeImage(capturedURL)
	if err != nil {
		return nil, fmt.Errorf("captured image download failed: %w", err)
	}

	req := (&facebody.CompareFaceAdvanceRequest{}).
		SetImageURLAObject(bytes.NewReader(imageA)).
		SetImageURLBObject(bytes.NewReader(imageB))
	runtime := (&dara.RuntimeOptions{}).
		SetConnectTimeout(compareConnectTimeout).
		SetReadTimeout(compareReadTimeout)
	return client.CompareFaceAdvance(req, runtime)
}

func downloadAndOptimizeImage(rawURL string) ([]byte, error) {
	if resolver, err := getMinioResolver(); err == nil && resolver != nil {
		if data, ok, err := resolver.fetch(rawURL); ok {
			if err != nil {
				log.Printf("[FaceAuth] minio fetch failed for %s: %v", rawURL, err)
			} else {
				optimized, optimizeErr := compressImageIfNeeded(data)
				if optimizeErr != nil {
					return data, nil
				}
				return optimized, nil
			}
		}
	}

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
		candidate := buf.Bytes()
		if len(candidate) < len(best) {
			best = append([]byte(nil), candidate...)
		}
		if len(best) <= targetImageBytes {
			break
		}
	}

	return best, nil
}

func getMinioResolver() (*minioResolver, error) {
	minioOnce.Do(func() {
		endpoint := strings.TrimSpace(os.Getenv("MINIO_ENDPOINT"))
		accessKey := strings.TrimSpace(os.Getenv("MINIO_ACCESS_KEY"))
		secretKey := strings.TrimSpace(os.Getenv("MINIO_SECRET_KEY"))
		bucket := strings.TrimSpace(os.Getenv("MINIO_BUCKET"))
		if endpoint == "" || accessKey == "" || secretKey == "" || bucket == "" {
			return
		}

		useSSL := strings.EqualFold(strings.TrimSpace(os.Getenv("MINIO_USE_SSL")), "true")
		client, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			sharedMinioErr = err
			return
		}

		sharedMinio = &minioResolver{
			client:    client,
			bucket:    bucket,
			publicURL: strings.TrimSuffix(strings.TrimSpace(os.Getenv("MINIO_PUBLIC_URL")), "/"),
		}
	})

	return sharedMinio, sharedMinioErr
}

func (r *minioResolver) fetch(rawURL string) ([]byte, bool, error) {
	objectKey := r.resolveObjectKey(rawURL)
	if objectKey == "" {
		return nil, false, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), downloadTimeout)
	defer cancel()

	object, err := r.client.GetObject(ctx, r.bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, true, err
	}
	defer object.Close()

	if _, err := object.Stat(); err != nil {
		return nil, true, err
	}

	reader := io.LimitReader(object, maxImageDownloadBytes+1)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, true, err
	}
	if len(data) == 0 {
		return nil, true, errors.New("empty image content")
	}
	if len(data) > maxImageDownloadBytes {
		return nil, true, fmt.Errorf("image is too large, max %d bytes", maxImageDownloadBytes)
	}
	return data, true, nil
}

func (r *minioResolver) resolveObjectKey(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}

	if r.publicURL != "" {
		publicPrefix := r.publicURL + "/" + r.bucket + "/"
		if strings.HasPrefix(rawURL, publicPrefix) {
			return strings.TrimPrefix(rawURL, publicPrefix)
		}
	}

	search := "/" + r.bucket + "/"
	if idx := strings.Index(rawURL, search); idx >= 0 {
		return strings.TrimPrefix(rawURL[idx+len(search):], "/")
	}

	return ""
}

func getCompareMode() string {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("ALIYUN_FACE_COMPARE_MODE")))
	switch mode {
	case "download", "auto", "url":
		return mode
	default:
		return "download"
	}
}
