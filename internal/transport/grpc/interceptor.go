package grpctransport

import (
	"context"
	"time"

	"github.com/RomanKovalev007/topqueries-service/internal/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryMetricsInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start).Seconds()

	st, _ := status.FromError(err)
	metrics.GRPCRequestsTotal.WithLabelValues(info.FullMethod, st.Code().String()).Inc()
	metrics.GRPCRequestDuration.WithLabelValues(info.FullMethod).Observe(duration)

	return resp, err
}
