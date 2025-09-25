package handlers

import (
	pkgLogger "notification/pkg/logger"
	servicesType "notification/pkg/services"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// universal handler for shipping services
// converter used in the services logic to process values

type KafkaMessageHandler[TRequest any, TService servicesType.Notifier[TRequest]] struct {
	Logger  pkgLogger.Logger
	Service TService
}

func NewKafkaMessageHandler[TRequest any, TService servicesType.Notifier[TRequest]](logger pkgLogger.Logger, service TService) *KafkaMessageHandler[TRequest, TService] {
	return &KafkaMessageHandler[TRequest, TService]{
		Logger:  logger,
		Service: service,
	}
}

func (h KafkaMessageHandler[TRequest, TService]) HandleMessage(msg *kafka.Message, converter func([]byte) (TRequest, error)) {
	h.Logger.Infof("Message from topic: " + *msg.TopicPartition.Topic)

	req, err := converter(msg.Value)
	if err != nil {
		h.Logger.Errorf("Failed parse message: %v", err)
		return
	}

	if err := h.Service.Send(req); err != nil {
		h.Logger.Errorf("Error sending msg: %v", err)
		return
	}

	h.Logger.Infof("Message send DONE")
}
