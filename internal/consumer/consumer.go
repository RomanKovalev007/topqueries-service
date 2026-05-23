package consumer

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/RomanKovalev007/topqueries-service/internal/domain"
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
	logReader *kafka.Reader
	svc       service
	log       *slog.Logger
}

func NewConsumer(cfg KafkaConfig, svc service, log *slog.Logger) *Consumer {
	return &Consumer{
		logReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: cfg.Brokers,
			Topic:   cfg.Topic,
			GroupID: cfg.GroupID,
		}),
		svc: svc,
		log: log,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	defer c.logReader.Close()
	for {
		var event domain.SearchEvent
		msg, err := c.logReader.ReadMessage(ctx)
		if err != nil{
			if ctx.Err() != nil{
				return
			}
			c.log.Error("failed to read kafka message", "err", err)
			continue
		}
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			c.log.Error("failed to unmarshal message", "err", err, "offset", msg.Offset)
			continue
		}
		if event.QueryText == ""{
			c.log.Warn("empty query, skipping", "offset", msg.Offset)
			continue
		}
		c.svc.Add(event)
	}
}