package svc

import (
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartStatsConsumer(svcCtx *ServiceContext) {
	if svcCtx.Config.RabbitMQ.URL() == "" {
		log.Println("RabbitMQ URL is empty in stats-rpc, consumer skipped.")
		return
	}

	// Wait a bit for connection
	time.Sleep(3 * time.Second)
	mqClient := svcCtx.MQ
	if mqClient == nil {
		log.Println("RabbitMQ client is nil in stats-rpc, consumer skipped.")
		return
	}

	err := mqClient.ConsumeEvents("ai_report.generate", func(delivery amqp.Delivery) {
		defer func() {
			_ = delivery.Ack(false)
		}()

		var task struct {
			ReportID int64 `json:"report_id"`
			UserID   int64 `json:"user_id"`
		}

		if err := json.Unmarshal(delivery.Body, &task); err != nil {
			log.Printf("[Stats Consumer] failed to unmarshal ai_report.generate task: %v", err)
			return
		}

		log.Printf("[Stats Consumer] received report generate task: %+v", task)

		err := svcCtx.ReportSvc.GenerateReportAsync(task.ReportID, task.UserID)
		if err != nil {
			log.Printf("[Stats Consumer] failed to generate report %d asynchronously: %v", task.ReportID, err)
		} else {
			log.Printf("[Stats Consumer] successfully generated report %d asynchronously", task.ReportID)
		}
	})

	if err != nil {
		log.Printf("failed to start ai_report.generate consumer: %v", err)
	} else {
		log.Println("Started ai_report.generate consumer successfully.")
	}
}
