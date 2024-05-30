package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/phonghaido/log-ingestor/helpers"
	"github.com/phonghaido/log-ingestor/types"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	KafkaConfig helpers.KafkaConfig
	Writer      kafka.Writer
	AckChan     map[string]chan bool
	AckMutex    sync.Mutex
}

func NewKafkaWriter(kafkaConfig helpers.KafkaConfig) kafka.Writer {
	log.Printf("Initializing Kafka writer for broker %s and topic %s", kafkaConfig.KafkaBroker, kafkaConfig.Topic)
	return *kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaConfig.KafkaBroker},
		Topic:   kafkaConfig.Topic,
	})
}

func NewKafkaProducer(kafkaConfig helpers.KafkaConfig) *KafkaProducer {
	return &KafkaProducer{
		KafkaConfig: kafkaConfig,
		Writer:      NewKafkaWriter(kafkaConfig),
		AckChan:     make(map[string]chan bool),
	}
}

func (p *KafkaProducer) CreateKafkaTopic() error {
	conn, err := kafka.Dial("tcp", p.KafkaConfig.KafkaBroker)
	if err != nil {
		log.Println("Error connecting to Kafka")
		return err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		log.Println("Error reading Kafka partitions")
		return err
	}
	for _, partition := range partitions {
		if partition.Topic == p.KafkaConfig.Topic {
			log.Printf("Topic %s already existed", p.KafkaConfig.Topic)
			return nil
		}
	}

	controller, err := conn.Controller()
	if err != nil {
		log.Println("Error creating Kafka controller")
		return err
	}

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Println("Error connecting to Kafka controller")
		return err
	}
	defer controllerConn.Close()

	topicConfig := []kafka.TopicConfig{
		{
			Topic:             p.KafkaConfig.Topic,
			NumPartitions:     p.KafkaConfig.NumPartition,
			ReplicationFactor: p.KafkaConfig.Replication,
		},
	}
	err = controllerConn.CreateTopics(topicConfig...)
	if err != nil {
		log.Printf("Error creating topic %s", p.KafkaConfig.Topic)
		return err
	}
	log.Printf("Topic %s created successfully", p.KafkaConfig.Topic)
	return nil
}

func (p *KafkaProducer) ProduceLogKafka(c echo.Context, logData types.LogData) error {
	logBytes, err := json.Marshal(logData)
	if err != nil {
		log.Println("Error marshalling log data")
		return err
	}

	p.AckMutex.Lock()
	p.AckChan[logData.ID] = make(chan bool)
	p.AckMutex.Unlock()

	err = p.Writer.WriteMessages(c.Request().Context(), kafka.Message{
		Value: logBytes,
	})
	log.Printf("Produce message: %s to kafka topic %s successfully.\n", logData.TraceID, p.KafkaConfig.Topic)
	if err != nil {
		return err
	}
	return err
}

func (p *KafkaProducer) WaitForAck(logID string, timeout time.Duration) (bool, error) {
	p.AckMutex.Lock()
	ackChan, exists := p.AckChan[logID]
	p.AckMutex.Unlock()
	if !exists {
		log.Printf("Channel for handling %s doesn't exist", logID)
		return false, fmt.Errorf("internal server error")
	}

	select {
	case <-ackChan:
		log.Printf("received acknowledgement for the log %s\n", logID)
		return true, nil
	case <-time.After(timeout):
		log.Println("timeout while waiting for acknowledgement")
		return false, fmt.Errorf("gateway timeout")
	}
}

func (p *KafkaProducer) Acknowledge(logID string) {
	p.AckMutex.Lock()
	if ackChan, exists := p.AckChan[logID]; exists {
		close(ackChan)
		delete(p.AckChan, logID)
	}
	p.AckMutex.Unlock()
}
