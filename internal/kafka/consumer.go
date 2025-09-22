package kafka

import (
	"context"
	json "encoding/json"
	models "notification/internal/models"
	services "notification/internal/services"
	logger "notification/pkg/logger"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type NotificationConsumer struct {
	broker   string
	groupID  string
	topics   []string
	service  services.EmailService
	producer Producer
	log      logger.Logger
}

type Producer interface {
	Publish(ctx context.Context, topic string, event models.EmailRequest) error
}

func NewNotificationConsumer(broker, groupID string, topics []string, svc services.EmailService, prod Producer, log logger.Logger) *NotificationConsumer {
	return &NotificationConsumer{
		broker:   broker,
		groupID:  groupID,
		topics:   topics,
		service:  svc,
		producer: prod,
		log:      log,
	}
}

func (c *NotificationConsumer) Start(ctx context.Context) {
	for _, topic := range c.topics {
		go c.consumeTopic(ctx, topic)
	}
}

func (c *NotificationConsumer) consumeTopic(ctx context.Context, topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{c.broker},
		GroupID:  c.groupID,
		Topic:    topic,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	c.log.Infof("Starting Kafka consumer for topic: %s", topic)

	for {
		select {
		case <-ctx.Done():
			c.log.Infof("Stopping consumer for topic: %s", topic)
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				c.log.Errorf("Error reading message: %v", err)
				time.Sleep(time.Second)
				continue
			}

			var emailReq models.EmailRequest
			if err := json.Unmarshal(msg.Value, &emailReq); err != nil {
				c.log.Errorf("Error unmarshaling Kafka message: %v", err)
				continue
			}

			if err := c.service.SendEmail(emailReq); err != nil {
				c.log.Errorf("Failed to send email to %s: %v", emailReq.To, err)
				_ = c.producer.Publish(ctx, "notification_failed", emailReq)
			} else {
				c.log.Infof("Email sent to %s", emailReq.To)
				_ = c.producer.Publish(ctx, "notification_send", emailReq)
			}
		}
	}
}
