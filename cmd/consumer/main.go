package main

import (
	"context"
	"log"

	"github.com/phonghaido/log-ingestor/db"
	"github.com/phonghaido/log-ingestor/helpers"
	"github.com/phonghaido/log-ingestor/kafka"
)

func main() {
	ctx := context.Background()
	kafkaConfig := helpers.ReadKafkaConfig()

	mongoClient, err := db.ConnectToMongoDB(ctx)
	if err != nil {
		log.Printf("Error connecting to mongodb: %s", err.Error())
	}

	logPersister := db.NewMongoDBLogPersister(mongoClient)
	kafkaConsumer := kafka.NewKafkaConsumer(*kafkaConfig, logPersister)
	kafkaConsumer.ConsumeLogKafka(ctx)
}
