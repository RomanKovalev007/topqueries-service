package consumer

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/RomanKovalev007/topqueries-service/internal/domain"
	"github.com/RomanKovalev007/topqueries-service/internal/metrics"
	"github.com/segmentio/kafka-go"
)

type service interface {
	Add(searchEvent domain.SearchEvent)
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type Consumer struct {
	reader *kafka.Reader
	svc       service
	log       *slog.Logger
}

func NewConsumer(cfg KafkaConfig, svc service, log *slog.Logger) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: cfg.Brokers,
			Topic:   cfg.Topic,
			GroupID: cfg.GroupID,
		}),
		svc: svc,
		log: log,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	defer c.reader.Close()
	c.log.Info("kafka consumer started", "topic", c.reader.Config().Topic, "brokers", c.reader.Config().Brokers)
	for {
		var event domain.SearchEvent
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			c.log.Error("failed to read kafka message", "err", err)
			metrics.KafkaMessagesTotal.WithLabelValues("read_error").Inc()
			continue
		}
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			c.log.Error("failed to unmarshal message", "err", err, "offset", msg.Offset)
			metrics.KafkaMessagesTotal.WithLabelValues("parse_error").Inc()
			continue
		}
		if event.QueryText == "" {
			c.log.Warn("empty query, skipping", "offset", msg.Offset)
			metrics.KafkaMessagesTotal.WithLabelValues("empty_query").Inc()
			continue
		}
		metrics.KafkaMessagesTotal.WithLabelValues("received").Inc()
		c.svc.Add(event)
	}
}