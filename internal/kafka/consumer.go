package kafka

import (
	"fmt"
	pkgLogger "notification/pkg/logger"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// func type for handlers
type MessageHandler func(msg *kafka.Message)

type Consumer struct {
	Consumer *kafka.Consumer
}

func NewConsumer(addrs []string, groupId string) (*Consumer, error) {
	c, err := kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers": strings.Join(addrs, ","),
			"group.id":          groupId,
			"auto.offset.reset": "earliest",
		},
	)

	if err != nil {
		return nil, fmt.Errorf("%s", "Error creating consumer: "+err.Error())
	}

	return &Consumer{Consumer: c}, nil
}

func (c *Consumer) Subscribe(topics []string) error {
	return c.Consumer.SubscribeTopics(topics, nil)
}

func (c *Consumer) Listen(handler map[string]MessageHandler, log pkgLogger.Logger) {
	go func() {
		for {
			ev := c.Consumer.Poll(2000)
			if ev == nil {
				log.Infof("Not Events")
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				m := e

				// result trasports ev
				if m.TopicPartition.Error != nil {
					log.Errorf("Delivery failed: %v\n", m.TopicPartition.Error)
					continue
				}

				if handler, ok := handler[*e.TopicPartition.Topic]; ok {
					handler(m)
				} else {
					log.Infof("Not found handler for this topic: " + *e.TopicPartition.Topic)
				}
			case *kafka.Error:
				log.Errorf("%s", "Error: "+e.Error())
			default:
				log.Infof("%s", "Ignored its event: "+ev.String())
			}
		}
	}()
}

func (c *Consumer) Close() {
	c.Consumer.Close()
}
