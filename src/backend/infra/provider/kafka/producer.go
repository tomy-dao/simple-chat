package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"local/config"
	"local/util/logger"

	"github.com/segmentio/kafka-go"
)

// Producer handles Kafka message production
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	logger.Info(nil, "Kafka producer initialized", map[string]interface{}{
		"brokers": brokers,
		"topic":   topic,
	})

	return &Producer{
		writer: writer,
	}
}

// ProduceMessage sends a message to Kafka
func (p *Producer) ProduceMessage(ctx context.Context, key string, value interface{}) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		logger.Error(nil, "Failed to marshal Kafka message", err)
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: valueBytes,
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		logger.Error(nil, "Failed to write message to Kafka", err)
		return fmt.Errorf("failed to write message: %w", err)
	}

	logger.Info(nil, "Message sent to Kafka", map[string]interface{}{
		"key":   key,
		"topic": p.writer.Topic,
	})

	return nil
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	if p.writer != nil {
		logger.Info(nil, "Closing Kafka producer", nil)
		return p.writer.Close()
	}
	return nil
}

// MessageProducer is a global producer for messages
var MessageProducer *Producer

// NotificationProducer is a global producer for notifications
var NotificationProducer *Producer

// InitProducers initializes all Kafka producers
func InitProducers() error {
	brokers := config.Config.KafkaBrokers
	if len(brokers) == 0 {
		return fmt.Errorf("kafka brokers not configured")
	}

	MessageProducer = NewProducer(brokers, config.Config.KafkaMessageTopic)
	NotificationProducer = NewProducer(brokers, config.Config.KafkaNotificationTopic)

	logger.Info(nil, "All Kafka producers initialized", map[string]interface{}{
		"brokers": brokers,
	})

	return nil
}

// CloseProducers closes all Kafka producers
func CloseProducers() {
	if MessageProducer != nil {
		MessageProducer.Close()
	}
	if NotificationProducer != nil {
		NotificationProducer.Close()
	}
}
