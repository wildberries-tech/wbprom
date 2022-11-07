package prom

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// httpServerMetrics is a struct that allows to write metrics of count and latency of http requests
type httpServerMetric struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

func NewHttpServerMetrics(appName string) *httpServerMetric {
	reqsCollector := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "reqs_count",
			Help:        "How many HTTP requests processed",
			ConstLabels: prometheus.Labels{"app": appName},
		},
		[]string{"method", "status", "path", "client"},
	)

	latencyCollector := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "reqs_latency",
		Help:        "How long it took to process the request",
		ConstLabels: prometheus.Labels{"app": appName},
		Buckets:     []float64{5, 10, 20, 30, 50, 70, 100, 150, 200, 300, 500, 1000},
	},
		[]string{"method", "status", "path", "client"},
	)

	prometheus.MustRegister(reqsCollector, latencyCollector)

	return &httpServerMetric{
		reqs:    reqsCollector,
		latency: latencyCollector,
	}
}

// Inc increases requests counter by one.
//  method, code, path and client are label values for "method", "status", "path" and "client" fields
func (h *httpServerMetric) Inc(method, code, path, client string) {
	h.reqs.WithLabelValues(method, code, path, client).Inc()
}

// WriteTiming writes time elapsed since the startTime.
// method, code, path and client are label values for "method", "status", "path" and "client" fields
func (h *httpServerMetric) WriteTiming(startTime time.Time, method, code, path, client string) {
	h.latency.WithLabelValues(method, code, path, client).Observe(timeFromStart(startTime))
}
