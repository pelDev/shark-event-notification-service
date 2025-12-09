package main

import (
	"context"
	"log"
	"net/smtp"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/commitshark/notification-svc/internal"
	"github.com/commitshark/notification-svc/internal/application/services"
	"github.com/commitshark/notification-svc/internal/domain/ports"
	grpcclient "github.com/commitshark/notification-svc/internal/infrastructure/adapters/grpc"
	"github.com/commitshark/notification-svc/internal/infrastructure/adapters/kafka"
	"github.com/commitshark/notification-svc/internal/infrastructure/adapters/providers"
	"github.com/commitshark/notification-svc/internal/infrastructure/adapters/sqlite"
	"github.com/commitshark/notification-svc/internal/infrastructure/adapters/templates"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	config := config.LoadConfig()

	// Initialize SMTP Client
	auth := smtp.PlainAuth("", config.Email.Username, config.Email.Password, config.Email.SMTPHost)

	// Initialize SQLite repository
	repo, err := sqlite.NewSQLiteNotificationRepository(config.SQLite.Path)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	// Initialize gRPC connection
	conn, err := grpc.NewClient(config.UserGrpcTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	renderer, err := templates.NewGoTemplateRenderer(templates.Files)
	if err != nil {
		log.Fatalf("failed to initialize template renderer: %v", err)
	}

	userDataAdapter := grpcclient.NewUserDataGRPCClient(conn)

	// Initialize providers
	providerList := []ports.NotificationProvider{
		providers.NewEmailProvider(
			config.Email.SMTPHost,
			config.Email.SMTPPort,
			config.Email.Username,
			config.Email.Password,
			config.Email.From,
			renderer,
			auth,
		),
		providers.NewSMSProvider(),
		providers.NewPushProvider(),
	}

	// Initialize service (no event publisher for now)
	notificationService := services.NewNotificationService(repo, providerList)

	// Initialize Kafka handler
	kafkaHandler := kafka.NewKafkaMessageHandler(notificationService, userDataAdapter)

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

func startRetryWorker(ctx context.Context, service *services.NotificationService, config struct {
	RetryBatchSize int           `mapstructure:"retry_batch_size"`
	RetryInterval  time.Duration `mapstructure:"retry_interval"`
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
