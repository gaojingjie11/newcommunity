package service

import (
	"context"
	"log"
	"time"

	"smartcommunity-microservices/app/stats/rpc/internal/model"

	"github.com/minio/minio-go/v7"
)

const (
	tempPayFacePrefix          = "face/temp-pay/"
	tempPayFaceCleanupInterval = 30 * time.Minute
	tempPayFaceTTL             = 30 * time.Minute
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
	go s.startTempPayFaceCleanupLoop()
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
	s.cleanupPrefix("", "garbage/", 0)
}

func (s *ReportScheduler) startTempPayFaceCleanupLoop() {
	if s.minioClient == nil || s.bucket == "" {
		log.Println("[Temp Face Cleanup] MinIO client or bucket not configured, periodic temp face cleanup skipped.")
		return
	}

	ticker := time.NewTicker(tempPayFaceCleanupInterval)
	defer ticker.Stop()

	log.Printf("[Temp Face Cleanup] Periodic cleanup started. Prefix=%s interval=%v ttl=%v", tempPayFacePrefix, tempPayFaceCleanupInterval, tempPayFaceTTL)

	for {
		select {
		case <-ticker.C:
			s.cleanupPrefix("[Temp Face Cleanup] ", tempPayFacePrefix, tempPayFaceTTL)
		case <-s.stopCh:
			log.Println("[Temp Face Cleanup] Stopping periodic temp face cleanup.")
			return
		}
	}
}

func (s *ReportScheduler) cleanupPrefix(logPrefix, prefix string, olderThan time.Duration) {
	if s.minioClient == nil || s.bucket == "" {
		log.Printf("%sMinIO client or bucket not configured, cleanup skipped.", logPrefix)
		return
	}

	if logPrefix == "" {
		logPrefix = "[Report Scheduler Cleanup] "
	}

	if olderThan > 0 {
		log.Printf("%sStarting cleanup for prefix '%s' older than %v...", logPrefix, prefix, olderThan)
	} else {
		log.Printf("%sStarting cleanup for all objects under prefix '%s'...", logPrefix, prefix)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	objectCh := s.minioClient.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	objectsToDelete := make(chan minio.ObjectInfo)
	queuedCount := 0
	cutoff := time.Now().Add(-olderThan)
	go func() {
		defer close(objectsToDelete)
		for object := range objectCh {
			if object.Err != nil {
				log.Printf("%serror listing object: %v", logPrefix, object.Err)
				continue
			}
			if olderThan > 0 && !object.LastModified.Before(cutoff) {
				continue
			}
			objectsToDelete <- object
			queuedCount++
		}
	}()

	errorCount := 0
	for rErr := range s.minioClient.RemoveObjects(ctx, s.bucket, objectsToDelete, minio.RemoveObjectsOptions{}) {
		log.Printf("%serror deleting object %s: %v", logPrefix, rErr.ObjectName, rErr.Err)
		errorCount++
	}

	successCount := queuedCount - errorCount
	log.Printf("%sCompleted. Successfully removed %d (out of %d) objects from MinIO bucket '%s' with prefix '%s'.", logPrefix, successCount, queuedCount, s.bucket, prefix)
}

func (s *ReportScheduler) Stop() {
	close(s.stopCh)
}
