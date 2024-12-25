package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

// NewKafkaConsumer initializes a new Kafka consumer
func NewKafkaConsumer(brokers []string, topic, groupID string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &KafkaConsumer{reader: reader}
}

// ConsumeMessages reads messages from the Kafka topic
func (c *KafkaConsumer) ConsumeMessages(ctx context.Context) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			break
		}
		log.Printf("Received message: Key=%s, Value=%s\n", string(m.Key), string(m.Value))
	}
}

// Close closes the Kafka consumer
func (c *KafkaConsumer) Close() {
	c.reader.Close()
}
