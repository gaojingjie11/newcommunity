package service

import (
	"context"
	"log"
	"time"

	"smartcommunity-microservices/app/stats/rpc/internal/model"

	"github.com/minio/minio-go/v7"
)

type ReportScheduler struct {
	reportSvc   *ReportService
	minioClient *minio.Client
	bucket      string
	stopCh      chan struct{}
}

func NewReportScheduler(reportSvc *ReportService, minioClient *minio.Client, bucket string) *ReportScheduler {
	return &ReportScheduler{
		reportSvc:   reportSvc,
		minioClient: minioClient,
		bucket:      bucket,
		stopCh:      make(chan struct{}),
	}
}

func (s *ReportScheduler) Start() {
	go func() {
		log.Println("[Report Scheduler] Starting automatic daily AI report scheduler...")
		for {
			now := time.Now()
			// Calculate the next occurrence of 9:00 PM (21:00:00)
			nextRun := time.Date(now.Year(), now.Month(), now.Day(), 21, 0, 0, 0, now.Location())
			if now.After(nextRun) {
				nextRun = nextRun.Add(24 * time.Hour)
			}
			duration := nextRun.Sub(now)
			log.Printf("[Report Scheduler] Next automatic report generation scheduled at %s (in %v)", nextRun.Format("2006-01-02 15:04:05"), duration)

			select {
			case <-time.After(duration):
				log.Println("[Report Scheduler] Daily trigger fired! Generating automatic AI report...")
				// operatorID = 0 represents system automatic generation
				var report *model.AIReport
				var err error
				report, err = s.reportSvc.GenerateReport(0)
				if err != nil {
					log.Printf("[Report Scheduler] Failed to generate automatic report: %v", err)
				} else {
					log.Printf("[Report Scheduler] Automatic report successfully generated with ID: %d", report.ID)
				}

				// Trigger MinIO garbage image cleanup
				log.Println("[Report Scheduler] Daily trigger fired! Triggering MinIO garbage image cleanup...")
				s.CleanupGarbage()

			case <-s.stopCh:
				log.Println("[Report Scheduler] Stopping automatic report scheduler.")
				return
			}
		}
	}()
}

func (s *ReportScheduler) CleanupGarbage() {
	if s.minioClient == nil || s.bucket == "" {
		log.Println("[Report Scheduler Cleanup] MinIO client or bucket not configured, skipping garbage image cleanup.")
		return
	}

	log.Println("[Report Scheduler Cleanup] Starting cleanup of garbage recognition images under 'garbage/' prefix...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 1. List objects with prefix "garbage/"
	objectCh := s.minioClient.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		Prefix:    "garbage/",
		Recursive: true,
	})

	// 2. Queue objects for deletion and count them
	objectsToDelete := make(chan minio.ObjectInfo)
	queuedCount := 0
	go func() {
		defer close(objectsToDelete)
		for object := range objectCh {
			if object.Err != nil {
				log.Printf("[Report Scheduler Cleanup] error listing object: %v", object.Err)
				continue
			}
			objectsToDelete <- object
			queuedCount++
		}
	}()

	// 3. RemoveObjects (returns errors only)
	errorCount := 0
	for rErr := range s.minioClient.RemoveObjects(ctx, s.bucket, objectsToDelete, minio.RemoveObjectsOptions{}) {
		log.Printf("[Report Scheduler Cleanup] error deleting object %s: %v", rErr.ObjectName, rErr.Err)
		errorCount++
	}

	successCount := queuedCount - errorCount
	log.Printf("[Report Scheduler Cleanup] Completed. Successfully removed %d (out of %d) garbage images from MinIO bucket '%s'.", successCount, queuedCount, s.bucket)
}

func (s *ReportScheduler) Stop() {
	close(s.stopCh)
}
