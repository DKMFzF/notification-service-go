package kafka

import (
	"context"
	"encoding/json"
	"notification/internal/models"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	broker string
}

func NewKafkaProducer(broker string) *KafkaProducer {
	return &KafkaProducer{broker: broker}
}

func (p *KafkaProducer) Publish(ctx context.Context, topic string, event models.EmailRequest) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{p.broker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return writer.WriteMessages(ctx, kafka.Message{Value: msg})
}
