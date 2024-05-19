package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"strconv"

	"github.com/phonghaido/log-ingestor/data"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	KafkaConfig data.KafkaConfig
	Writer      kafka.Writer
}

func NewKafkaProducer(kafkaConfig data.KafkaConfig) *KafkaProducer {
	return &KafkaProducer{
		KafkaConfig: kafkaConfig,
		Writer: *kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{kafkaConfig.KafkaBroker},
			Topic:   kafkaConfig.Topic,
		}),
	}
}

func (p *KafkaProducer) ConnectToKafka() (*kafka.Conn, error) {
	conn, err := kafka.Dial("tcp", p.KafkaConfig.KafkaBroker)
	if err != nil {
		log.Println("Error connecting to Kafka")
		return nil, err
	}

	return conn, nil
}

func (p *KafkaProducer) CheckTopicExist(conn *kafka.Conn) (bool, error) {
	partitions, err := conn.ReadPartitions()
	if err != nil {
		log.Println("Error reading Kafka partitions")
		return false, err
	}
	for _, partition := range partitions {
		if partition.Topic == p.KafkaConfig.Topic {
			return true, nil
		}
	}
	return false, nil
}

func (p *KafkaProducer) CreateKafkaTopic() error {
	conn, err := p.ConnectToKafka()
	if err != nil {
		return err
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

func (p *KafkaProducer) ProduceLogKafka(logData data.LogData) error {
	defer p.Writer.Close()

	logBytes, err := json.Marshal(logData)
	if err != nil {
		log.Println("Error marshalling log data")
		return err
	}

	err = p.Writer.WriteMessages(context.Background(), kafka.Message{
		Value: logBytes,
	})
	log.Println(err.Error())
	return err
}
