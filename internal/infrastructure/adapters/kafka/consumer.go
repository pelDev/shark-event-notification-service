package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/commitshark/notification-svc/internal"
	"github.com/commitshark/notification-svc/internal/domain/events"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader        *kafka.Reader
	handler       *KafkaMessageHandler
	topic         string
	consumerGroup string
	brokers       []string
	logger        *log.Logger
}

func NewKafkaConsumer(
	kConfig config.KafkaConfig,
	handler *KafkaMessageHandler,
) *KafkaConsumer {
	logger := log.New(os.Stdout, "[KafkaConsumer] ", log.LstdFlags)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        kConfig.Brokers,
		Topic:          kConfig.Topic,
		GroupID:        kConfig.ConsumerGroup,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		MaxWait:        time.Second,
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
		// Logger:         kafka.LoggerFunc(logger.Printf),
		ErrorLogger: kafka.LoggerFunc(logger.Printf),
	})

	return &KafkaConsumer{
		reader:        reader,
		handler:       handler,
		topic:         kConfig.Topic,
		consumerGroup: kConfig.ConsumerGroup,
		brokers:       kConfig.Brokers,
		logger:        logger,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) error {
	c.logger.Printf("Starting Kafka consumer for topic: %s", c.topic)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		c.logger.Println("Shutdown signal received")
		c.reader.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			c.logger.Println("Context cancelled, stopping consumer")
			return c.reader.Close()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					return nil
				}
				c.logger.Printf("Error fetching message: %v", err)
				time.Sleep(time.Second)
				continue
			}

			// Process message
			if err := c.processMessage(ctx, msg); err != nil {
				c.logger.Printf("Failed to process message: %v", err)
				// TODO: send to DLQ
				continue
			}

			// Commit offset
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				c.logger.Printf("Failed to commit message: %v", err)
			}
		}
	}
}

func (c *KafkaConsumer) processMessage(ctx context.Context, msg kafka.Message) error {
	var request events.KafkaEvent

	if err := json.Unmarshal(msg.Value, &request); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	c.logger.Printf("Processing event: %s ---> Type %s", request.ID, request.EventType)

	// Validate
	if err := request.Validate(); err != nil {
		return fmt.Errorf("failed to validate request: %w", err)
	}

	return c.handler.HandleMessage(ctx, msg.Value)
}

func (c *KafkaConsumer) Close() error {
	c.logger.Println("Closing Kafka consumer")
	return c.reader.Close()
}
