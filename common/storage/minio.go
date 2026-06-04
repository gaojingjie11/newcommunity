package storage

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOConfig struct {
	Endpoint  string `json:",optional"`
	AccessKey string `json:",optional"`
	SecretKey string `json:",optional"`
	Bucket    string `json:",optional"`
	UseSSL    bool   `json:",optional"`
}

func Init(cfg MinIOConfig) (*minio.Client, error) {
	return minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
}
