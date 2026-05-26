package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCPort    string `env:"GRPC_PORT"    env-default:":50051"`
	MetricsPort string `env:"METRICS_PORT" env-default:":9090"`
	LogLevel    string `env:"LOG_LEVEL"    env-default:"info"`
	Kafka KafkaConfig
}

type KafkaConfig struct {
	Brokers string `env:"KAFKA_BROKERS" env-default:"localhost:9092"`
	Topic   string `env:"KAFKA_TOPIC" env-default:"search-events"`
	GroupID string `env:"KAFKA_GROUP_ID" env-default:"topqueries-service"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return &cfg, nil
}