package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/phonghaido/log-ingestor/data"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	KafkaConfig data.KafkaConfig
	Writer      kafka.Writer
	AckChan     map[string]chan bool
	AckMutex    sync.Mutex
}

func NewKafkaProducer(kafkaConfig data.KafkaConfig) *KafkaProducer {
	return &KafkaProducer{
		KafkaConfig: kafkaConfig,
		Writer:      data.NewKafkaWriter(kafkaConfig),
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
		log.Fatalf("Error creating topic %s", p.KafkaConfig.Topic)
		return err
	}
	log.Printf("Topic %s created successfully", p.KafkaConfig.Topic)
	return nil
}

func (p *KafkaProducer) ProduceLogKafka(logData data.LogData) (string, error) {
	logBytes, err := json.Marshal(logData)
	if err != nil {
		log.Println("Error marshalling log data")
		return "", err
	}

	p.AckMutex.Lock()
	p.AckChan[logData.TraceID] = make(chan bool)
	p.AckMutex.Unlock()

	err = p.Writer.WriteMessages(context.Background(), kafka.Message{
		Value: logBytes,
	})
	log.Printf("Produce message: %s to kafka topic %s successfully.\n", logData.TraceID, p.KafkaConfig.Topic)
	if err != nil {
		return "", err
	}
	return logData.TraceID, err
}

func (p *KafkaProducer) WaitForAck(traceID string, timeout time.Duration) (bool, error) {
	p.AckMutex.Lock()
	ackChan, exists := p.AckChan[traceID]
	p.AckMutex.Unlock()
	if !exists {
		log.Printf("Channel for handling %s doesn't exist", traceID)
		return false, fmt.Errorf("internal server error")
	}
	select {
	case <-ackChan:
		log.Printf("received acknowledgement for the log %s\n", traceID)
		return true, nil
	case <-time.After(timeout):
		log.Println("timeout while waiting for acknowledgement")
		return false, fmt.Errorf("gateway timeout")
	}
}

func (p *KafkaProducer) Acknowledge(traceID string) {
	p.AckMutex.Lock()
	if ackChan, exists := p.AckChan[traceID]; exists {
		close(ackChan)
		delete(p.AckChan, traceID)
	}
	p.AckMutex.Unlock()
}
