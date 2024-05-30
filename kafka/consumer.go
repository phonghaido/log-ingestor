package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/phonghaido/log-ingestor/helpers"
	"github.com/phonghaido/log-ingestor/types"
	"github.com/segmentio/kafka-go"
)

type LogPersister interface {
	PersistLog(context.Context, types.LogData) error
}

type KafkaConsumer struct {
	KafkaConfig helpers.KafkaConfig
	Reader      kafka.Reader
	Persister   LogPersister
}

func NewKafkaReader(kafkaConfig helpers.KafkaConfig) kafka.Reader {
	return *kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaConfig.KafkaBroker},
		Topic:   kafkaConfig.Topic,
		GroupID: "log-processor",
	})
}

func NewKafkaConsumer(kafkaConfig helpers.KafkaConfig, logPersister LogPersister) *KafkaConsumer {
	return &KafkaConsumer{
		KafkaConfig: kafkaConfig,
		Reader:      NewKafkaReader(kafkaConfig),
		Persister:   logPersister,
	}
}

func (c *KafkaConsumer) ConsumeLogKafka(ctx context.Context) {
	defer c.Reader.Close()
	for {
		msg, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Failed to read message: %s\n", err.Error())
			continue
		}
		log.Printf("Received message: %s\n", string(msg.Value))
		var logData types.LogData
		err = json.Unmarshal(msg.Value, &logData)
		if err != nil {
			log.Printf("Error decoding message %s\n", err.Error())
			continue
		}

		if err := c.Persister.PersistLog(ctx, logData); err != nil {
			continue
		}

		log.Printf("Insert record %s to mongodb collection successfully", logData.TraceID)

		err = SendAcknowledgement(logData.ID)
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
