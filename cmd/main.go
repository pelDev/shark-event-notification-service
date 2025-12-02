package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/commitshark/notification-svc/internal/application/services"
	"github.com/commitshark/notification-svc/internal/infrastructure/adapters/kafka"
	"github.com/commitshark/notification-svc/internal/infrastructure/adapters/providers"
	"github.com/commitshark/notification-svc/internal/infrastructure/adapters/sqlite"
	// "github.com/segmentio/kafka-go/protocol/consumer"
)

type Config struct {
	SQLite struct {
		Path string
	}
	Kafka struct {
		Brokers       []string
		Topic         string
		ConsumerGroup string
	}
	Email struct {
		SMTPHost string
		SMTPPort int
		Username string
		Password string
		From     string
	}
	Service struct {
		RetryBatchSize int
		RetryInterval  time.Duration
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	config := loadConfig()

	// Initialize SQLite repository
	repo, err := sqlite.NewSQLiteNotificationRepository(config.SQLite.Path)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	// Initialize provider (simplified for now)
	provider := providers.NewEmailProvider(
		config.Email.SMTPHost,
		config.Email.SMTPPort,
		config.Email.Username,
		config.Email.Password,
		config.Email.From,
	)

	// Initialize service (no event publisher for now)
	notificationService := services.NewNotificationService(repo, provider)

	// Initialize Kafka handler
	kafkaHandler := kafka.NewKafkaMessageHandler(notificationService)

	// Initialize Kafka consumer (using segmentio/kafka-go)
	consumer := kafka.NewKafkaConsumer(config.Kafka, kafkaHandler)

	// Start retry worker
	go startRetryWorker(ctx, notificationService, config.Service)

	// Start Kafka consumer
	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("Kafka consumer error: %v", err)
			cancel()
		}
	}()

	// Graceful shutdown
	waitForShutdown(cancel)

	log.Println("Notification service shutdown complete")
}

func loadConfig() Config {
	// In practice, use viper, envconfig, or similar
	return Config{
		SQLite: struct{ Path string }{
			Path: "./data/notifications.db",
		},
		Kafka: struct {
			Brokers       []string
			Topic         string
			ConsumerGroup string
		}{
			Brokers:       []string{"localhost:9092"},
			Topic:         "notifications",
			ConsumerGroup: "notification-service",
		},
		Email: struct {
			SMTPHost string
			SMTPPort int
			Username string
			Password string
			From     string
		}{
			SMTPHost: os.Getenv("SMTP_HOST"),
			SMTPPort: 587,
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     os.Getenv("SMTP_FROM"),
		},
		Service: struct {
			RetryBatchSize int
			RetryInterval  time.Duration
		}{
			RetryBatchSize: 100,
			RetryInterval:  30 * time.Second,
		},
	}
}

func startRetryWorker(ctx context.Context, service *services.NotificationService, config struct {
	RetryBatchSize int
	RetryInterval  time.Duration
}) {
	ticker := time.NewTicker(config.RetryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := service.RetryFailedNotifications(ctx, config.RetryBatchSize); err != nil {
				log.Printf("Retry worker error: %v", err)
			}
		}
	}
}

func waitForShutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received signal: %v", sig)
	cancel()

	// Give services time to shutdown gracefully
	time.Sleep(2 * time.Second)
}
