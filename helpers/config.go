package helpers

import (
	"log"
	"os"
	"strconv"
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

func ReadKafkaConfig() *KafkaConfig {
	partition, err := strconv.Atoi(os.Getenv("NUM_PARTITION"))
	if err != nil {
		log.Printf("Error reading env variable")
	}
	replication, err := strconv.Atoi(os.Getenv("REPLICATION"))
	if err != nil {
		log.Printf("Error reading env variable")
	}
	return NewKafkaConfig(
		os.Getenv("KAFKA_BROKER"),
		os.Getenv("KAFKA_TOPIC"),
		partition,
		replication,
	)
}
