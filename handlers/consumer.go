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
		var logData data.LogData
		err = json.Unmarshal(msg.Value, &logData)
		if err != nil {
			log.Printf("Error decoding message %s\n", err.Error())
			continue
		}

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
