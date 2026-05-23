package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/RomanKovalev007/topqueries-service/gen/pb"
	"github.com/RomanKovalev007/topqueries-service/internal/config"
	"github.com/RomanKovalev007/topqueries-service/internal/consumer"
	"github.com/RomanKovalev007/topqueries-service/internal/service"
	"github.com/RomanKovalev007/topqueries-service/internal/store"
	grpctransport "github.com/RomanKovalev007/topqueries-service/internal/transport/grpc"
	"google.golang.org/grpc"
)


func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	window := store.NewWindow()
	go window.Ticker(ctx)

	svc := service.NewService(window)

	kafkaCfg := consumer.KafkaConfig{
		Brokers: strings.Split(cfg.Kafka.Brokers, ","),
		Topic: cfg.Kafka.Topic,
		GroupID: cfg.Kafka.GroupID,
	}
	c := consumer.NewConsumer(kafkaCfg, svc)
	go c.Start(ctx)

	srv := grpctransport.NewServer(svc)

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTopQueriesServiceServer(grpcServer, srv)


	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("grpc server error: %v", err)
		}
	}()

	<-ctx.Done()
	grpcServer.GracefulStop()
}
