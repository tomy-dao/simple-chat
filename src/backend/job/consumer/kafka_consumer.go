package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"local/config"
	"local/util/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

// MessageHandler is a function that handles consumed messages
type MessageHandler func(ctx context.Context, key []byte, value []byte) error

// KafkaConsumer manages Kafka message consumption
type KafkaConsumer struct {
	reader  *kafka.Reader
	handler MessageHandler
	topic   string
}

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(brokers []string, groupID string, topic string, handler MessageHandler) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	logger.Info(nil, "Kafka consumer initialized", map[string]interface{}{
		"brokers": brokers,
		"group":   groupID,
		"topic":   topic,
	})

	return &KafkaConsumer{
		reader:  reader,
		handler: handler,
		topic:   topic,
	}
}

// Start starts consuming messages and blocks until interrupted
func (kc *KafkaConsumer) Start() error {
	logger.Info(nil, "Starting Kafka consumer", map[string]interface{}{
		"topic": kc.topic,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start consuming in a goroutine
	errChan := make(chan error, 1)
	go func() {
		for {
			msg, err := kc.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					// Context cancelled, normal shutdown
					errChan <- nil
					return
				}
				logger.Error(nil, "Error reading Kafka message", err)
				errChan <- fmt.Errorf("failed to read message: %w", err)
				return
			}

			logger.Info(nil, "Received Kafka message", map[string]interface{}{
				"topic":     msg.Topic,
				"partition": msg.Partition,
				"offset":    msg.Offset,
			})

			// Handle the message
			if err := kc.handler(ctx, msg.Key, msg.Value); err != nil {
				logger.Error(nil, "Error handling Kafka message", err)
				// Continue processing other messages even if one fails
			}
		}
	}()

	// Wait for shutdown signal or error
	select {
	case <-sigChan:
		logger.Info(nil, "Received shutdown signal, stopping consumer", nil)
		cancel()
		return <-errChan
	case err := <-errChan:
		return err
	}
}

// Stop closes the Kafka consumer
func (kc *KafkaConsumer) Stop() error {
	if kc.reader != nil {
		logger.Info(nil, "Closing Kafka consumer", map[string]interface{}{
			"topic": kc.topic,
		})
		return kc.reader.Close()
	}
	return nil
}

// MessageEvent represents a chat message event
type MessageEvent struct {
	MessageID      uint   `json:"message_id"`
	ConversationID uint   `json:"conversation_id"`
	UserID         uint   `json:"user_id"`
	Content        string `json:"content"`
}

// NotificationEvent represents a notification event
type NotificationEvent struct {
	UserID  uint   `json:"user_id"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

// MessageEventHandler handles message events from Kafka
func MessageEventHandler(ctx context.Context, key []byte, value []byte) error {
	var event MessageEvent
	if err := json.Unmarshal(value, &event); err != nil {
		logger.Error(nil, "Failed to unmarshal message event", err)
		return fmt.Errorf("failed to unmarshal message event: %w", err)
	}

	logger.Info(nil, "Processing message event", map[string]interface{}{
		"message_id":      event.MessageID,
		"conversation_id": event.ConversationID,
		"user_id":         event.UserID,
	})

	// TODO: Add your business logic here
	// For example: update read status, send notifications, etc.

	return nil
}

// NotificationEventHandler handles notification events from Kafka
func NotificationEventHandler(ctx context.Context, key []byte, value []byte) error {
	var event NotificationEvent
	if err := json.Unmarshal(value, &event); err != nil {
		logger.Error(nil, "Failed to unmarshal notification event", err)
		return fmt.Errorf("failed to unmarshal notification event: %w", err)
	}

	logger.Info(nil, "Processing notification event", map[string]interface{}{
		"user_id": event.UserID,
		"type":    event.Type,
		"message": event.Message,
	})

	// TODO: Add your business logic here
	// For example: send push notification, email, etc.

	return nil
}

// StartMessageConsumer starts the message consumer
func StartMessageConsumer() error {
	consumer := NewKafkaConsumer(
		config.Config.KafkaBrokers,
		config.Config.KafkaConsumerGroup,
		config.Config.KafkaMessageTopic,
		MessageEventHandler,
	)
	defer consumer.Stop()
	return consumer.Start()
}

// StartNotificationConsumer starts the notification consumer
func StartNotificationConsumer() error {
	consumer := NewKafkaConsumer(
		config.Config.KafkaBrokers,
		config.Config.KafkaConsumerGroup,
		config.Config.KafkaNotificationTopic,
		NotificationEventHandler,
	)
	defer consumer.Stop()
	return consumer.Start()
}
