package consumer

import (
	"context"
	"encoding/json"
	"log"

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
}

func NewConsumer(cfg KafkaConfig, svc service) *Consumer {
	return &Consumer{
		logReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: cfg.Brokers,
			Topic:   cfg.Topic,
			GroupID: cfg.GroupID,
		}),
		svc: svc,
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
			log.Printf("failed to read log message: %v", err)
			continue
		}
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to unmarshal log message: %v", err)
			continue
		}
		if event.QueryText == ""{
			log.Print("empty query, skipping")
			continue
		}
		c.svc.Add(event)
	}
}