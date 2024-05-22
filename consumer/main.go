package main

import (
	"github.com/phonghaido/log-ingestor/handlers"
	"github.com/phonghaido/log-ingestor/helpers"
)

func main() {
	kafkaConfig := helpers.ReadKafkaConfig()
	kafkaConsumer := handlers.NewKafkaConsumer(*kafkaConfig)
	kafkaConsumer.ConsumeLogKafka()
}
