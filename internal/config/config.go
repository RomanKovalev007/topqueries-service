package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCPort    string `env:"GRPC_PORT"    env-default:":50051"`
	MetricsPort string `env:"METRICS_PORT" env-default:":9090"`
	LogLevel    string `env:"LOG_LEVEL"    env-default:"info"`
	Kafka       KafkaConfig
	Window      WindowConfig
	RateLimiter RateLimiterConfig
}

type KafkaConfig struct {
	Brokers string `env:"KAFKA_BROKERS"   env-default:"localhost:9092"`
	Topic   string `env:"KAFKA_TOPIC"     env-default:"search-events"`
	GroupID string `env:"KAFKA_GROUP_ID"  env-default:"topqueries-service"`
}

type WindowConfig struct {
	Duration    time.Duration `env:"WINDOW_DURATION"     env-default:"5m"`
	BucketCount int           `env:"WINDOW_BUCKET_COUNT" env-default:"30"`
	MaxTopN     int           `env:"WINDOW_MAX_TOP_N"    env-default:"100"`
}

type RateLimiterConfig struct {
	WindowDuration  time.Duration `env:"RATE_LIMITER_WINDOW_DURATION"   env-default:"1m"`
	BucketCount int           `env:"RATE_LIMITER_BUCKET_COUNT"  env-default:"12"`
	MaxRequests int64         `env:"RATE_LIMITER_MAX_REQUESTS"  env-default:"100"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return &cfg, nil
}
