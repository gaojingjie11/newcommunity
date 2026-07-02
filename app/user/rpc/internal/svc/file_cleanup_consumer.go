package svc

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
)

type FileCleanupEvent struct {
	URL string `json:"url"`
}

func StartFileCleanupConsumer(svcCtx *ServiceContext) {
	if svcCtx.MqClient == nil {
		log.Println("RabbitMQ client is nil, file cleanup consumer skipped.")
		return
	}
	if svcCtx.MinioClient == nil {
		log.Println("MinIO client is nil, file cleanup consumer skipped.")
		return
	}

	// Wait a bit for RabbitMQ connection to stabilize
	time.Sleep(3 * time.Second)

	err := svcCtx.MqClient.ConsumeEvents("file.cleanup", func(delivery amqp.Delivery) {
		var event FileCleanupEvent
		if err := json.Unmarshal(delivery.Body, &event); err != nil {
			log.Printf("[File Cleanup Consumer] failed to unmarshal file.cleanup event: %v", err)
			_ = delivery.Reject(false)
			return
		}

		url := strings.TrimSpace(event.URL)
		if url == "" {
			_ = delivery.Ack(false)
			return
		}

		log.Printf("[File Cleanup Consumer] received file.cleanup event for url: %s", url)

		bucketName := svcCtx.Config.MinIO.Bucket
		// Parse object key from url
		// Expected url format: http://minio:9000/smartcommunity/image/12345.png
		searchStr := "/" + bucketName + "/"
		idx := strings.Index(url, searchStr)
		var objectKey string
		if idx != -1 {
			objectKey = url[idx+len(searchStr):]
		} else {
			// Fallback: split by slash and join from index 4 onwards
			parts := strings.Split(url, "/")
			if len(parts) >= 5 {
				objectKey = strings.Join(parts[4:], "/")
			}
		}

		if objectKey == "" {
			log.Printf("[File Cleanup Consumer] failed to resolve object key for url: %s", url)
			_ = delivery.Reject(false)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Printf("[File Cleanup Consumer] deleting object from MinIO: bucket=%s, key=%s", bucketName, objectKey)
		err := svcCtx.MinioClient.RemoveObject(ctx, bucketName, objectKey, minio.RemoveObjectOptions{})
		if err != nil {
			log.Printf("[File Cleanup Consumer] failed to remove object from MinIO: %v", err)
			_ = delivery.Nack(false, true)
		} else {
			log.Printf("[File Cleanup Consumer] successfully removed object from MinIO: bucket=%s, key=%s", bucketName, objectKey)
			_ = delivery.Ack(false)
		}
	})

	if err != nil {
		log.Printf("failed to start file.cleanup consumer: %v", err)
	} else {
		log.Println("Started file.cleanup consumer successfully.")
	}
}
