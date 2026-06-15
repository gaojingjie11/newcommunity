// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"smartcommunity-microservices/app/gateway/api/internal/logic"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/common/response"

	miniogo "github.com/minio/minio-go/v7"
)

func UploadGarbageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
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
		dir := "garbage"
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

		protocol := "http"
		if cfg.UseSSL {
			protocol = "https"
		}
		fileURL := fmt.Sprintf("%s://%s/%s/%s", protocol, cfg.Endpoint, cfg.Bucket, info.Key)

		l := logic.NewUploadGarbageLogic(r.Context(), svcCtx)
		resp, err := l.UploadGarbage(fileURL, header.Filename)
		response.Response(w, resp, err)
	}
}
