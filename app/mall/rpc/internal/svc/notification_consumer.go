package svc

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"smartcommunity-microservices/app/mall/rpc/internal/service"
	"smartcommunity-microservices/common/mail"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartNotificationConsumer(svcCtx *ServiceContext) {
	if svcCtx.Config.RabbitMQ.URL() == "" {
		log.Println("RabbitMQ URL is empty, notification consumer skipped.")
		return
	}

	// Wait a bit for RabbitMQ client to connect
	time.Sleep(3 * time.Second)
	mqClient := svcCtx.MQ
	if mqClient == nil {
		log.Println("RabbitMQ client is nil, notification consumer skipped.")
		return
	}

	// 1. Consume order.paid
	err := mqClient.ConsumeEvents(service.QueueOrderPaid, func(delivery amqp.Delivery) {
		defer func() {
			_ = delivery.Ack(false)
		}()

		var event service.OrderPaidEvent
		if err := json.Unmarshal(delivery.Body, &event); err != nil {
			log.Printf("[Notification Consumer] failed to unmarshal order.paid event: %v", err)
			return
		}

		log.Printf("[Notification Consumer] received order.paid event: %+v", event)

		// Fetch user profile from database to get email
		user, err := svcCtx.UserRepo.FindByID(event.UserID)
		if err != nil {
			log.Printf("[Notification Consumer] failed to find user %d: %v", event.UserID, err)
			return
		}

		if user.Email == "" {
			log.Printf("[Notification Consumer] user %d (%s) does not have email bound, skipping notification.", event.UserID, user.Username)
			return
		}

		// Fetch order details
		order, err := svcCtx.OrderRepo.FindByID(event.OrderID)
		if err != nil {
			log.Printf("[Notification Consumer] failed to find order %d: %v", event.OrderID, err)
			return
		}

		// Fetch payment record to identify specific payment method
		paymentRecord, err := svcCtx.PaymentRepo.FindByOrderID(event.OrderID)
		paymentMethodText := "余额支付"
		if err == nil && paymentRecord != nil {
			switch paymentRecord.PaymentMethod {
			case "alipay":
				paymentMethodText = "支付宝支付"
			case "face":
				paymentMethodText = "人脸支付"
			case "password":
				paymentMethodText = "密码支付"
			case "nopassword":
				paymentMethodText = "免密支付"
			case "balance", "wallet":
				paymentMethodText = "余额支付"
			default:
				paymentMethodText = "电子支付"
			}
		}

		// Map items
		var items []mail.OrderItemInfo
		for _, item := range order.Items {
			productName := item.ProductSnapshot
			if productName == "" && item.Product.Name != "" {
				productName = item.Product.Name
			}
			if productName == "" {
				productName = "未知商品"
			}
			items = append(items, mail.OrderItemInfo{
				ProductName: productName,
				PriceCents:  item.Price,
				Quantity:    item.Quantity,
			})
		}

		paidAtText := time.Now().Format("2006-01-02 15:04:05")
		if order.PaidAt != nil {
			paidAtText = order.PaidAt.Format("2006-01-02 15:04:05")
		} else if event.PaidAt != "" {
			if t, err := time.Parse(time.RFC3339, event.PaidAt); err == nil {
				paidAtText = t.Format("2006-01-02 15:04:05")
			}
		}

		// Send email!
		username := user.RealName
		if username == "" {
			username = user.Username
		}

		subject := fmt.Sprintf("【智能社区】订单 %s 支付成功通知", order.OrderNo)
		bodyHtml := mail.GenerateOrderPaidHTML(username, order.OrderNo, order.TotalAmount, order.UsedPoints, order.UsedBalance, paymentMethodText, paidAtText, items)

		if err := mail.SendMail(svcCtx.Config.Mail, user.Email, subject, bodyHtml); err != nil {
			log.Printf("[Notification Consumer] failed to send order.paid email to %s: %v", user.Email, err)
		} else {
			log.Printf("[Notification Consumer] successfully sent order.paid email to %s", user.Email)
		}
	})
	if err != nil {
		log.Printf("failed to start order.paid consumer: %v", err)
	} else {
		log.Println("Started order.paid consumer successfully.")
	}

	// 2. Consume wallet.recharged
	err = mqClient.ConsumeEvents(service.QueueWalletRecharged, func(delivery amqp.Delivery) {
		defer func() {
			_ = delivery.Ack(false)
		}()

		var event service.WalletRechargedEvent
		if err := json.Unmarshal(delivery.Body, &event); err != nil {
			log.Printf("[Notification Consumer] failed to unmarshal wallet.recharged event: %v", err)
			return
		}

		log.Printf("[Notification Consumer] received wallet.recharged event: %+v", event)

		// Fetch user profile from database
		user, err := svcCtx.UserRepo.FindByID(event.UserID)
		if err != nil {
			log.Printf("[Notification Consumer] failed to find user %d: %v", event.UserID, err)
			return
		}

		if user.Email == "" {
			log.Printf("[Notification Consumer] user %d (%s) does not have email bound, skipping notification.", event.UserID, user.Username)
			return
		}

		paymentMethodText := "模拟支付"
		if strings.HasPrefix(event.IdempotencyKey, "RECH_") {
			paymentMethodText = "支付宝支付"
		}

		rechargedAtText := time.Now().Format("2006-01-02 15:04:05")
		if t, err := time.Parse(time.RFC3339, event.RechargedAt); err == nil {
			rechargedAtText = t.Format("2006-01-02 15:04:05")
		}

		// Send email!
		username := user.RealName
		if username == "" {
			username = user.Username
		}

		subject := "【智能社区】电子钱包充值成功通知"
		bodyHtml := mail.GenerateWalletRechargedHTML(username, event.Amount, paymentMethodText, rechargedAtText, event.IdempotencyKey)

		if err := mail.SendMail(svcCtx.Config.Mail, user.Email, subject, bodyHtml); err != nil {
			log.Printf("[Notification Consumer] failed to send wallet.recharged email to %s: %v", user.Email, err)
		} else {
			log.Printf("[Notification Consumer] successfully sent wallet.recharged email to %s", user.Email)
		}
	})
	if err != nil {
		log.Printf("failed to start wallet.recharged consumer: %v", err)
	} else {
		log.Println("Started wallet.recharged consumer successfully.")
	}

	// 3. Consume order timeout trigger
	err = mqClient.ConsumeEvents(service.QueueOrderTimeout, func(delivery amqp.Delivery) {
		defer func() {
			_ = delivery.Ack(false)
		}()

		var event service.OrderCancelledEvent
		if err := json.Unmarshal(delivery.Body, &event); err != nil {
			log.Printf("[Notification Consumer] failed to unmarshal order.timeout event: %v", err)
			return
		}

		log.Printf("[Notification Consumer] received order.timeout event: %+v", event)

		if err := svcCtx.TimeoutSvc.CancelExpiredOrder(event.OrderID); err != nil {
			log.Printf("[Notification Consumer] failed to cancel expired order %d: %v", event.OrderID, err)
		} else {
			log.Printf("[Notification Consumer] successfully cancelled expired order %d", event.OrderID)
		}
	})
	if err != nil {
		log.Printf("failed to start order.timeout consumer: %v", err)
	} else {
		log.Println("Started order.timeout consumer successfully.")
	}

	// 4. Consume property_fee.paid
	err = mqClient.ConsumeEvents("property_fee.paid", func(delivery amqp.Delivery) {
		defer func() {
			_ = delivery.Ack(false)
		}()

		var event struct {
			Event          string `json:"event"`
			FeeID          int64  `json:"fee_id"`
			UserID         int64  `json:"user_id"`
			Month          string `json:"month"`
			Amount         int64  `json:"amount"`
			PaidAt         string `json:"paid_at"`
			IdempotencyKey string `json:"idempotency_key"`
		}
		if err := json.Unmarshal(delivery.Body, &event); err != nil {
			log.Printf("[Notification Consumer] failed to unmarshal property_fee.paid event: %v", err)
			return
		}

		log.Printf("[Notification Consumer] received property_fee.paid event: %+v", event)

		// Fetch user profile from database to get email
		user, err := svcCtx.UserRepo.FindByID(event.UserID)
		if err != nil {
			log.Printf("[Notification Consumer] failed to find user %d: %v", event.UserID, err)
			return
		}

		if user.Email == "" {
			log.Printf("[Notification Consumer] user %d (%s) does not have email bound, skipping property fee notification.", event.UserID, user.Username)
			return
		}

		paidAtText := time.Now().Format("2006-01-02 15:04:05")
		if t, err := time.Parse(time.RFC3339, event.PaidAt); err == nil {
			paidAtText = t.Format("2006-01-02 15:04:05")
		}

		username := user.RealName
		if username == "" {
			username = user.Username
		}

		subject := fmt.Sprintf("【智能社区】物业缴费成功通知 - %s", event.Month)
		bodyHtml := mail.GeneratePropertyFeePaidHTML(username, event.Month, event.Amount, "电子钱包支付", paidAtText, event.IdempotencyKey)

		if err := mail.SendMail(svcCtx.Config.Mail, user.Email, subject, bodyHtml); err != nil {
			log.Printf("[Notification Consumer] failed to send property_fee.paid email to %s: %v", user.Email, err)
		} else {
			log.Printf("[Notification Consumer] successfully sent property_fee.paid email to %s", user.Email)
		}
	})
	if err != nil {
		log.Printf("failed to start property_fee.paid consumer: %v", err)
	} else {
		log.Println("Started property_fee.paid consumer successfully.")
	}
}
