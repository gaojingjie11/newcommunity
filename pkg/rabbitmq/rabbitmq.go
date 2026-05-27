package rabbitmq

import (
	"context"
	"time"

	"smartcommunity-microservices/pkg/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func Connect(cfg config.RabbitMQConfig) (*Client, error) {
	conn, err := amqp.DialConfig(cfg.URL(), amqp.Config{Dial: amqp.DefaultDial(3 * time.Second)})
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return &Client{conn: conn, ch: ch}, nil
}

func (c *Client) Close() error {
	if c == nil {
		return nil
	}
	if c.ch != nil {
		_ = c.ch.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) PublishEvent(ctx context.Context, queue string, body []byte) error {
	if _, err := c.ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		return err
	}
	return c.ch.PublishWithContext(ctx, "", queue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Body:         body,
	})
}

func (c *Client) ConsumeEvents(queue string, handler func(amqp.Delivery)) error {
	if _, err := c.ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		return err
	}
	deliveries, err := c.ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for delivery := range deliveries {
			handler(delivery)
		}
	}()
	return nil
}
