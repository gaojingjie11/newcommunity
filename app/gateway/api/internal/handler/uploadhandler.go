package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/common/response"

	miniogo "github.com/minio/minio-go/v7"
)

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if svcCtx.MinioClient == nil {
			response.Response(w, nil, errors.New("对象存储不可用"))
			return
		}

		// Maximum upload size 10MB
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			response.Response(w, nil, fmt.Errorf("文件大小超出限制或解析错误: %w", err))
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			response.Response(w, nil, errors.New("请选择要上传的文件"))
			return
		}
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext == "" {
			ext = ".bin"
		}
		dir := uploadDir(r.PostFormValue("dir"), header.Filename, header.Header.Get("Content-Type"))
		objectName := fmt.Sprintf("%s/%d%s", dir, time.Now().UnixNano(), ext)
		contentType := header.Header.Get("Content-Type")

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		cfg := svcCtx.Config.MinIO
		info, err := svcCtx.MinioClient.PutObject(ctx, cfg.Bucket, objectName, file, header.Size, miniogo.PutObjectOptions{ContentType: contentType})
		if err != nil {
			response.Response(w, nil, fmt.Errorf("文件上传失败: %w", err))
			return
		}

		publicURL := os.Getenv("MINIO_PUBLIC_URL")
		var fileURL string
		if publicURL != "" {
			fileURL = fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(publicURL, "/"), cfg.Bucket, info.Key)
		} else {
			protocol := "http"
			if cfg.UseSSL {
				protocol = "https"
			}
			fileURL = fmt.Sprintf("%s://%s/%s/%s", protocol, cfg.Endpoint, cfg.Bucket, info.Key)
		}

		response.Response(w, &types.UploadResp{
			Url: fileURL,
			Key: info.Key,
		}, nil)
	}
}

func uploadDir(explicit, filename, contentType string) string {
	explicit = strings.Trim(strings.ToLower(explicit), "/ ")
	if explicit == "face/temp-pay" {
		return explicit
	}
	switch explicit {
	case "face", "image", "common":
		return explicit
	}
	name := strings.ToLower(filename)
	if strings.Contains(name, "face") {
		return "face"
	}
	if strings.HasPrefix(contentType, "image/") {
		return "image"
	}
	return "common"
}
