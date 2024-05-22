package main

import (
	"github.com/phonghaido/log-ingestor/data"
	"github.com/phonghaido/log-ingestor/handlers"
)

var (
	kafkaConfig   = data.NewKafkaConfig("localhost:9092", "logs", 1, 1)
	kafkaConsumer = handlers.NewKafkaConsumer(*kafkaConfig)
)

func main() {
	kafkaConsumer.ConsumeLogKafka()
}
