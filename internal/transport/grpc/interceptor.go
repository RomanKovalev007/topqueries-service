package grpctransport

import (
	"context"
	"log/slog"
	"time"

	"github.com/RomanKovalev007/topqueries-service/internal/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		st, _ := status.FromError(err)
		code := st.Code().String()

		metrics.GRPCRequestsTotal.WithLabelValues(info.FullMethod, code).Inc()
		metrics.GRPCRequestDuration.WithLabelValues(info.FullMethod).Observe(duration.Seconds())

		log.Info("grpc request",
			slog.String("method", info.FullMethod),
			slog.String("status", code),
			slog.Duration("duration", duration),
		)

		return resp, err
	}
}
