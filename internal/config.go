package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type EmailSMTPConfig struct {
	SMTPHost string `mapstructure:"smtp_host"`
	SMTPPort int    `mapstructure:"smtp_port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

type HttpEmailConfig struct {
	Url    string `mapstructure:"url"`
	APIKey string `mapstructure:"api_key"`
}

type ServiceConfig struct {
	RetryBatchSize int           `mapstructure:"retry_batch_size"`
	RetryInterval  time.Duration `mapstructure:"retry_interval"`
}

type KafkaConfig struct {
	Brokers       []string `mapstructure:"brokers"`
	Topic         string   `mapstructure:"topic"`
	ConsumerGroup string   `mapstructure:"consumer_group"`
}

type SQLiteConfig struct {
	Path string `mapstructure:"path"`
}

type Config struct {
	SQLite         SQLiteConfig    `mapstructure:"sqlite"`
	Kafka          KafkaConfig     `mapstructure:"kafka"`
	Email          EmailSMTPConfig `mapstructure:"email"`
	MarketingEmail EmailSMTPConfig `mapstructure:"marketing_email"`
	HTTPEmail      HttpEmailConfig `mapstructure:"http_email"`
	Service        ServiceConfig   `mapstructure:"service"`
	UserGrpcTarget string          `mapstructure:"user_grpc_target"`
}

func LoadConfig() Config {
	viper.AddConfigPath(".") // look in current directory
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("/root")

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Transactional email
	_ = viper.BindEnv("email.smtp_host", "EMAIL_SMTP_HOST")
	_ = viper.BindEnv("email.smtp_port", "EMAIL_SMTP_PORT")
	_ = viper.BindEnv("email.username", "EMAIL_USERNAME")
	_ = viper.BindEnv("email.password", "EMAIL_PASSWORD")
	_ = viper.BindEnv("email.from", "EMAIL_FROM")

	// Marketing email
	_ = viper.BindEnv("marketing_email.smtp_host", "MARKETING_EMAIL_SMTP_HOST")
	_ = viper.BindEnv("marketing_email.smtp_port", "MARKETING_EMAIL_SMTP_PORT")
	_ = viper.BindEnv("marketing_email.username", "MARKETING_EMAIL_USERNAME")
	_ = viper.BindEnv("marketing_email.password", "MARKETING_EMAIL_PASSWORD")
	_ = viper.BindEnv("marketing_email.from", "MARKETING_EMAIL_FROM")

	// HTTP email
	_ = viper.BindEnv("http_email.api_key", "HTTP_EMAIL_API_KEY")

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Loaded config file:", viper.ConfigFileUsed())
	} else {
		log.Println("No config file found, using environment variables")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	return cfg
}
