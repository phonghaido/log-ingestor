package handlers

import (
	"context"
	"log"

	"github.com/phonghaido/log-ingestor/data"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	KafkaConfig data.KafkaConfig
	Reader      kafka.Reader
}

func NewKafkaConsumer(kafkaConfig data.KafkaConfig) *KafkaConsumer {
	return &KafkaConsumer{
		KafkaConfig: kafkaConfig,
		Reader:      data.NewKafkaReader(kafkaConfig),
	}
}

func (c *KafkaConsumer) ConsumeLogKafka() {
	defer c.Reader.Close()

	for {
		msg, err := c.Reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Failed to read message: %s\n", err.Error())
			continue
		}
		log.Printf("Received message: %s\n", string(msg.Value))
	}
}
