package data

import "github.com/segmentio/kafka-go"

type KafkaConfig struct {
	KafkaBroker  string
	Topic        string
	NumPartition int
	Replication  int
}

func NewKafkaConfig(kafkaBroker, topic string, numPartition, replication int) *KafkaConfig {
	return &KafkaConfig{
		KafkaBroker:  kafkaBroker,
		Topic:        topic,
		NumPartition: numPartition,
		Replication:  replication,
	}
}

func NewKafkaWriter(kafkaConfig KafkaConfig) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaConfig.KafkaBroker},
		Topic:   kafkaConfig.Topic,
	})
}
