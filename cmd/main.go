package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/RomanKovalev007/topqueries-service/gen/pb"
	"github.com/RomanKovalev007/topqueries-service/internal/config"
	"github.com/RomanKovalev007/topqueries-service/internal/consumer"
	"github.com/RomanKovalev007/topqueries-service/internal/metrics"
	"github.com/RomanKovalev007/topqueries-service/internal/service"
	"github.com/RomanKovalev007/topqueries-service/internal/store"
	grpctransport "github.com/RomanKovalev007/topqueries-service/internal/transport/grpc"
	"github.com/RomanKovalev007/topqueries-service/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	l := logger.New(cfg.LogLevel)

	metrics.Register()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	window := store.NewQueryWindow(cfg.Window.Duration, cfg.Window.BucketCount, cfg.Window.MaxTopN)
	go window.Ticker(ctx)

	rateLimiter := store.NewRateLimiter(cfg.RateLimiter.WindowDuration, cfg.RateLimiter.BucketCount, cfg.RateLimiter.MaxRequests)
	go rateLimiter.Ticker(ctx)

	stoplist := store.NewStopList([]string{})

	svc := service.NewService(window, stoplist, rateLimiter, cfg.Window.Duration, l)

	kafkaCfg := consumer.KafkaConfig{
		Brokers: strings.Split(cfg.Kafka.Brokers, ","),
		Topic:   cfg.Kafka.Topic,
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

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpctransport.UnaryInterceptor(l)),
	)
	pb.RegisterTopQueriesServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			l.Error("grpc server stopped", "err", err)
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(cfg.MetricsPort, nil); err != nil {
			l.Error("metrics server stopped", "err", err)
		}
	}()

	l.Info("service started",
		slog.String("grpc_port", cfg.GRPCPort),
		slog.String("metrics_port", cfg.MetricsPort),
		slog.String("kafka_brokers", cfg.Kafka.Brokers),
		slog.String("kafka_topic", cfg.Kafka.Topic),
		slog.String("window_duration", cfg.Window.Duration.String()),
		slog.Int("window_buckets", cfg.Window.BucketCount),
		slog.Int("window_max_top_n", cfg.Window.MaxTopN),
		slog.String("rate_limiter_window", cfg.RateLimiter.WindowDuration.String()),
		slog.Int64("rate_limiter_max_requests", cfg.RateLimiter.MaxRequests),
	)

	<-ctx.Done()
	l.Info("shutting down")
	grpcServer.GracefulStop()
}
