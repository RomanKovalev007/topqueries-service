package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/RomanKovalev007/topqueries-service/gen/pb"
	"github.com/RomanKovalev007/topqueries-service/internal/consumer"
	"github.com/RomanKovalev007/topqueries-service/internal/service"
	"github.com/RomanKovalev007/topqueries-service/internal/store"
	grpctransport "github.com/RomanKovalev007/topqueries-service/internal/transport/grpc"
	"google.golang.org/grpc"
)


func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	window := store.NewWindow()
	go window.Ticker(ctx)
	
	svc := service.NewService(window)

	kafkaCfg := consumer.KafkaConfig{
		Brokers: []string{"localhost:9092"},
		Topic: "search-events",
		GroupID: "topqueries-service",
	}
	c := consumer.NewConsumer(kafkaCfg, svc)
	go c.Start(ctx)

	srv := grpctransport.NewServer(svc)

	lis, err := net.Listen("tcp", ":50051")
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
