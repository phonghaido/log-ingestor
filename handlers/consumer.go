package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/phonghaido/log-ingestor/data"
	"github.com/phonghaido/log-ingestor/db"
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
	ctx := context.Background()
	for {
		msg, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Failed to read message: %s\n", err.Error())
			continue
		}
		log.Printf("Received message: %s\n", string(msg.Value))
		var logData data.LogData
		err = json.Unmarshal(msg.Value, &logData)
		if err != nil {
			log.Printf("Error decoding message %s\n", err.Error())
			continue
		}

		mongoClient, err := db.ConnectToMongoDB(ctx)
		if err != nil {
			log.Printf("Error connecting to mongodb %s", err.Error())
			continue
		}

		err = db.InsertToDB(ctx, mongoClient, logData)
		if err != nil {
			log.Printf("Error inserting record to mongodb collection %s", err.Error())
			continue
		}

		log.Printf("Insert record %s to mongodb collection successfully", logData.TraceID)

		err = SendAcknowledgement(logData.TraceID)
		if err != nil {
			log.Printf("Error sending acknowledgement: %s", err.Error())
			continue
		}
	}
}

func SendAcknowledgement(logID string) error {
	url := os.Getenv("PRODUCER") + "/ack"
	ackData := map[string]string{"logId": logID}

	ackBytes, err := json.Marshal(ackData)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(ackBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send acknowledgement, status code: %d", resp.StatusCode)
	}
	return nil
}
