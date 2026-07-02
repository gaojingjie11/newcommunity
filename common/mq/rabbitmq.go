package mq

import (
	"context"
	"fmt"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	Host     string `json:",optional"`
	Port     int    `json:",optional"`
	Username string `json:",optional"`
	Password string `json:",optional"`
}

func (c RabbitMQConfig) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", c.Username, c.Password, c.Host, c.Port)
}

type Client struct {
	conn *amqp.Connection
}

func Connect(cfg RabbitMQConfig) (*Client, error) {
	conn, err := amqp.DialConfig(cfg.URL(), amqp.Config{Dial: amqp.DefaultDial(3 * time.Second)})
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	_ = ch.Close()
	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	if c == nil {
		return nil
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) PublishEvent(ctx context.Context, queue string, body []byte) error {
	if c == nil || c.conn == nil {
		return nil
	}
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if _, err := ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		return err
	}
	return ch.PublishWithContext(ctx, "", queue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Body:         body,
	})
}

func (c *Client) PublishDelayEvent(ctx context.Context, delayQueue string, targetQueue string, delayMs int64, body []byte) error {
	if c == nil || c.conn == nil {
		return nil
	}
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	args := amqp.Table{
		"x-dead-letter-exchange":    "", // default exchange
		"x-dead-letter-routing-key": targetQueue,
	}
	if _, err := ch.QueueDeclare(delayQueue, true, false, false, false, args); err != nil {
		return err
	}
	if _, err := ch.QueueDeclare(targetQueue, true, false, false, false, nil); err != nil {
		return err
	}
	if delayMs < 0 {
		delayMs = 0
	}
	return ch.PublishWithContext(ctx, "", delayQueue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Expiration:   strconv.FormatInt(delayMs, 10),
		Timestamp:    time.Now(),
		Body:         body,
	})
}

func (c *Client) ConsumeEvents(queue string, handler func(amqp.Delivery)) error {
	if c == nil || c.conn == nil {
		return nil
	}
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	if _, err := ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		_ = ch.Close()
		return err
	}
	deliveries, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		return err
	}
	go func() {
		for delivery := range deliveries {
			handler(delivery)
		}
		_ = ch.Close()
	}()
	return nil
}
