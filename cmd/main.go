package main

import (
	"context"
	"log"
	"log/slog"
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
	"github.com/RomanKovalev007/topqueries-service/pkg/logger"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	l := logger.New(cfg.LogLevel)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	window := store.NewWindow()
	go window.Ticker(ctx)

	stoplist := store.NewStopList([]string{})

	svc := service.NewService(window, stoplist)

	kafkaCfg := consumer.KafkaConfig{
		Brokers: strings.Split(cfg.Kafka.Brokers, ","),
		Topic: cfg.Kafka.Topic,
		GroupID: cfg.Kafka.GroupID,
	}
	c := consumer.NewConsumer(kafkaCfg, svc, l)
	go c.Start(ctx)

	srv := grpctransport.NewServer(svc)

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		l.Error("failed to listen", "err", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTopQueriesServiceServer(grpcServer, srv)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			l.Error("grpc server stopped", "err", err)
		}
	}()

	l.Info("service started", slog.String("grpc_port", cfg.GRPCPort))

	<-ctx.Done()
	l.Info("shutting down")
	grpcServer.GracefulStop()
}
