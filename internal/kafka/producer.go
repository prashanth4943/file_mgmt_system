package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer initializes a new Kafka producer
func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}
	return &KafkaProducer{writer: writer}
}

// SendMessage sends a message to the Kafka topic
func (p *KafkaProducer) SendMessage(key, value string) error {
	err := p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(value),
		},
	)
	if err != nil {
		log.Printf("Failed to send message: %v\n", err)
	}
	return err
}

// Close closes the Kafka producer
func (p *KafkaProducer) Close() {
	p.writer.Close()
}
