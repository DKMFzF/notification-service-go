package kafka

import (
	"errors"
	"fmt"
	"strings"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	flushTimeout = 5000 // 5s
)

var errUnknownType = errors.New("unknown event type")

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(addres []string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		// list brokers...
		"bootstrap.servers": strings.Join(addres, ","),
	}

	p, err := kafka.NewProducer(conf)

	if err != nil {
		return nil, fmt.Errorf("%s", "Error with create producer: "+err.Error())
	}

	return &Producer{producer: p}, nil
}

// send message in topic
func (p *Producer) Produce(message, topic string) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
		Key:   nil,
	}

	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return nil
	}

	e := <-kafkaChan

	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case *kafka.Error:
		return ev
	default:
		return errUnknownType
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
