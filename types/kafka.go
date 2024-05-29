package types

import (
	"log"

	"github.com/segmentio/kafka-go"
)

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

func NewKafkaWriter(kafkaConfig KafkaConfig) kafka.Writer {
	log.Printf("Initializing Kafka writer for broker %s and topic %s", kafkaConfig.KafkaBroker, kafkaConfig.Topic)
	return *kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaConfig.KafkaBroker},
		Topic:   kafkaConfig.Topic,
	})
}

func NewKafkaReader(kafkaConfig KafkaConfig) kafka.Reader {
	return *kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaConfig.KafkaBroker},
		Topic:   kafkaConfig.Topic,
		GroupID: "log-processor",
	})
}
