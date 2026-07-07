package faceauth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	facebody "github.com/alibabacloud-go/facebody-20191230/v6/client"
	"github.com/alibabacloud-go/tea/dara"
)

const (
	defaultEndpoint         = "facebody.cn-shanghai.aliyuncs.com"
	defaultRegion           = "cn-shanghai"
	defaultMinConfidence    = float32(75)
	defaultQualityThreshold = float32(70)
)

type compareConfig struct {
	endpoint         string
	region           string
	minConfidence    float32
	qualityThreshold float32
}

type CompareResult struct {
	Confidence    float64 `json:"confidence"`
	MinConfidence float64 `json:"min_confidence"`
	QualityA      float64 `json:"quality_a"`
	QualityB      float64 `json:"quality_b"`
	MessageTips   string  `json:"message_tips"`
}

var (
	loadOnce     sync.Once
	sharedClient *facebody.Client
	sharedCfg    compareConfig
	sharedErr    error
)

func ValidateEnrollment(ctx context.Context, imageURL string) error {
	_ = ctx
	imageURL = strings.TrimSpace(imageURL)
	if imageURL == "" {
		return errors.New("请先上传人脸照片")
	}

	client, _, err := getClient()
	if err != nil {
		return err
	}

	req := (&facebody.DetectFaceRequest{}).
		SetImageURL(imageURL).
		SetQuality(true).
		SetMaxFaceNumber(2)

	resp, err := client.DetectFaceWithOptions(req, &dara.RuntimeOptions{})
	if err != nil {
		return errors.New("人脸注册校验失败，请稍后重试")
	}
	if resp == nil || resp.Body == nil || resp.Body.Data == nil || resp.Body.Data.FaceCount == nil {
		return errors.New("人脸注册校验失败，请重新上传清晰正脸照片")
	}

	switch *resp.Body.Data.FaceCount {
	case 0:
		return errors.New("未识别到人脸，请上传清晰正脸照片")
	case 1:
		return nil
	default:
		return errors.New("检测到多张人脸，请上传仅包含本人的照片")
	}
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

	req := (&facebody.CompareFaceRequest{}).
		SetImageURLA(registeredURL).
		SetImageURLB(capturedURL).
		SetQualityScoreThreshold(cfg.qualityThreshold)

	resp, err := client.CompareFaceWithOptions(req, &dara.RuntimeOptions{})
	if err != nil {
		return nil, errors.New("人脸比对失败，请稍后重试")
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

	if result.MessageTips != "" {
		return result, errors.New(formatMessageTips(result.MessageTips))
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
