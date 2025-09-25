package kafka

type KafkaTopics struct {
	Topics map[string]map[string]KafkaEvent `yaml:"topics"`
}

type KafkaEvent struct {
	Service   string `yaml:"service"`
	Converter string `yaml:"converter"`
}
