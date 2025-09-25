package config

import (
	"os"
	"strconv"
	"strings"

	kafkaType "notification/pkg/kafka"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	godotenv "github.com/joho/godotenv"
	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Port              string
	SMTPUser          string
	SMTPPass          string
	SMTPHost          string
	SMTPPort          string
	KafkaBroker       string
	KafkaGroupId      string
	KafkaListenUpdate int
	KafkaTopics       kafkaType.KafkaTopics
}

func Load() *Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	listenTimer, _ := strconv.Atoi(os.Getenv("KAFKA_LISTEN_UPDATE"))
	if listenTimer == 0 {
		listenTimer = 2000
	}

	// read config for kafka
	kafkaFile, err := os.ReadFile("consumer.kafka.conf.yml")
	if err != nil {
		panic(err)
	}
	var kafkaTopics kafkaType.KafkaTopics
	if err := yaml.Unmarshal(kafkaFile, &kafkaTopics); err != nil {
		panic(err)
	}

	return &Config{
		Port:              port,
		SMTPUser:          os.Getenv("SMTP_USER"),
		SMTPPass:          os.Getenv("SMTP_PASS"),
		SMTPHost:          os.Getenv("SMTP_HOST"),
		SMTPPort:          os.Getenv("SMTP_PORT"),
		KafkaBroker:       os.Getenv("KAFKA_BROKER"),
		KafkaGroupId:      os.Getenv("KAFKA_GROUP_ID"),
		KafkaListenUpdate: listenTimer,
		KafkaTopics:       kafkaTopics,
	}
}

func ProducerConfig(addrs []string, groupId string) (*kafka.Consumer, error) {
	return kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers":        strings.Join(addrs, ","),
			"broker.address.family":    "v4",
			"group.id":                 groupId,
			"auto.offset.reset":        "earliest",
			"session.timeout.ms":       6000,
			"enable.auto.offset.store": false,
		},
	)
}
