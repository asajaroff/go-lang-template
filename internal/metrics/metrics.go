// Package metrics exposes Prometheus/VictoriaMetrics-format metrics.
//
// VictoriaMetrics is wire-compatible with Prometheus's scrape protocol, so
// nothing here is Prometheus- vs. VictoriaMetrics-specific — this is a
// natural complement to Factor XI (Logs): logs are the event stream,
// metrics are the aggregated counters/timings, but both are collected by
// the runtime environment, not buffered by the app itself.
package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests processed, labeled by method, route, and status.",
	}, []string{"method", "route", "status"})

	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency in seconds, labeled by method and route.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route"})
)

// Excluded from the metrics above: liveness/readiness are polled by
// Kubernetes every few seconds and /metrics scraping isn't application
// traffic — counting either would just add high-frequency noise.
var excludedRoutes = map[string]bool{
	"/healthz/liveness":  true,
	"/healthz/readiness": true,
	"/metrics":           true,
}

// Middleware records request count and latency for every request served by
// an http.ServeMux, labeled by the matched route *pattern* (e.g.
// "GET /items/{id}"), not the raw URL — using the raw URL would create
// unbounded label cardinality (one series per item ID).
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		route := r.Pattern
		if route == "" {
			route = r.URL.Path
		}
		if excludedRoutes[route] {
			return
		}

		requestsTotal.WithLabelValues(r.Method, route, strconv.Itoa(rec.status)).Inc()
		requestDuration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
	})
}

// Handler serves the Prometheus text exposition format at whatever path
// it's mounted on (see cmd/server, mounted at /metrics).
func Handler() http.Handler {
	return promhttp.Handler()
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
