package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/phonghaido/log-ingestor/db"
	"github.com/phonghaido/log-ingestor/helpers"
	"github.com/phonghaido/log-ingestor/kafka"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()
	kafkaConfig := helpers.ReadKafkaConfig()

	mongoClient, err := db.ConnectToMongoDB(ctx)
	if err != nil {
		log.Printf("Error connecting to mongodb: %s", err.Error())
	}

	logPersister := db.NewMongoDBLogPersister(mongoClient)

	kafkaConsumer := kafka.NewKafkaConsumer(*kafkaConfig, logPersister)
	log.Println("Waiting for messages...")

	// for partition := range kafkaConfig.NumPartition {

	// }
	kafkaConsumer.ConsumeLogKafka(ctx)
}
