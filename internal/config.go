package config

import (
	"log"
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
	SQLite  SQLiteConfig    `mapstructure:"sqlite"`
	Kafka   KafkaConfig     `mapstructure:"kafka"`
	Email   EmailSMTPConfig `mapstructure:"email"`
	Service ServiceConfig   `mapstructure:"service"`
}

func LoadConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // look in current directory
	viper.AutomaticEnv()     // override with env variables if present

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	return cfg
}
