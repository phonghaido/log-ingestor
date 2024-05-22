package helpers

import (
	"log"
	"os"
	"strconv"

	"github.com/phonghaido/log-ingestor/data"
)

func ReadKafkaConfig() *data.KafkaConfig {
	partition, err := strconv.Atoi(os.Getenv("NUM_PARTITION"))
	if err != nil {
		log.Printf("Error reading env variable")
	}
	replication, err := strconv.Atoi(os.Getenv("REPLICATION"))
	if err != nil {
		log.Printf("Error reading env variable")
	}
	return data.NewKafkaConfig(
		os.Getenv("KAFKA_BROKER"),
		os.Getenv("KAFKA_TOPIC"),
		partition,
		replication,
	)
}
