package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OrdersCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "crypto_sync_bot_orders_total",
		Help: "The total number of processed orders",
	}, []string{"exchange", "status"})

	OrderLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "crypto_sync_bot_order_latency_seconds",
		Help:    "Latency of order execution in seconds",
		Buckets: prometheus.DefBuckets,
	})

	StreamLag = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "crypto_sync_bot_stream_lag_seconds",
		Help: "Lag of the data stream in seconds",
	})
)
