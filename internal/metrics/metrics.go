package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	EventsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topqueries_events_total",
			Help: "Total number of incoming search events by result",
		},
		[]string{"result"},
	)

	GRPCRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topqueries_grpc_requests_total",
			Help: "Total number of gRPC requests by method and status",
		},
		[]string{"method", "status"},
	)

	GRPCRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topqueries_grpc_request_duration_seconds",
			Help:    "gRPC request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	KafkaMessagesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topqueries_kafka_messages_total",
			Help: "Total number of Kafka messages by result",
		},
		[]string{"result"},
	)
)

func Register() {
	prometheus.MustRegister(EventsTotal, GRPCRequestsTotal, GRPCRequestDuration, KafkaMessagesTotal)
}
