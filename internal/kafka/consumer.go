package kafka

import (
	"context"
	"fmt"
	"notification/internal/config"
	pkgLogger "notification/pkg/logger"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// func type for handlers
type MessageHandler func(msg *kafka.Message)

type Consumer struct {
	Consumer *kafka.Consumer
}

func NewConsumer(addrs []string, groupId string) (*Consumer, error) {
	c, err := config.ProducerConfig(addrs, groupId)

	if err != nil {
		return nil, fmt.Errorf("%s", "Error creating consumer: "+err.Error())
	}

	return &Consumer{Consumer: c}, nil
}

// subscribe on kafka topic
func (c *Consumer) Subscribe(topics []string) error {
	return c.Consumer.SubscribeTopics(topics, nil)
}

// listener kafka events
func (c *Consumer) Listen(ctx context.Context, handlers map[string]map[string]MessageHandler, log pkgLogger.Logger, timeListenUpdate int) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				continue
			default:
				ev := c.Consumer.Poll(timeListenUpdate)
				if ev == nil {
					log.Debugf("Not Events")
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

					topic := m.TopicPartition.Topic
					key := string(m.Key)

					// checks topics and event in cnf services
					if topicMap, ok := handlers[*topic]; ok {
						if handler, ok := topicMap[key]; ok {
							handler(m)
						} else {
							log.Infof("No found event %s in topic %s", key, topic)
						}
					} else {
						log.Infof("No found handler for this topic: " + *topic)
					}

				case *kafka.Error:
					log.Errorf("%s", "Error: "+e.Error())
				default:
					log.Warnf("%s", "Ignored its event: "+ev.String())
				}
			}
		}
	}()
}

func (c *Consumer) Close() {
	c.Consumer.Close()
}
