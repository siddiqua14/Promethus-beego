package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const timeoutThreshold = 5.0 // Seconds for timeout threshold

var (
	// Total outgoing HTTP requests
	httpOutgoingRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_outgoing_requests_total",
			Help: "Total number of outgoing HTTP requests",
		},
		[]string{"method", "url", "status"},
	)

	// Histogram for request duration
	httpOutgoingRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_outgoing_request_duration_histogram_seconds",
			Help:    "Histogram of outgoing HTTP request latencies",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "url", "status"},
	)

	// Track requests that exceed the timeout threshold
	httpOutgoingRequestTimeoutTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_outgoing_request_timeout_total",
			Help: "Total number of HTTP requests that exceeded timeout threshold",
		},
		[]string{"method", "url"},
	)

	initOnce sync.Once
)

// InstrumentedHttpClient is the struct for our HTTP client with monitoring
type InstrumentedHttpClient struct {
	Client *http.Client
}

// NewInstrumentedHttpClient creates a new instrumented HTTP client with the given timeout
func NewInstrumentedHttpClient(timeout time.Duration) *InstrumentedHttpClient {
	return &InstrumentedHttpClient{
		Client: &http.Client{Timeout: timeout},
	}
}

// Do is a custom method for making HTTP requests and recording metrics
func (c *InstrumentedHttpClient) Do(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := c.Client.Do(req)
	duration := time.Since(start).Seconds()

	// Default to "error" if request fails
	statusCode := "error"
	if err == nil {
		statusCode = http.StatusText(resp.StatusCode)
	}

	// Record Prometheus metrics for the request
	httpOutgoingRequestsTotal.WithLabelValues(req.Method, req.URL.Path, statusCode).Inc()
	httpOutgoingRequestDuration.WithLabelValues(req.Method, req.URL.Path, statusCode).Observe(duration)

	// If duration exceeds the timeout threshold, track it
	if duration > timeoutThreshold {
		httpOutgoingRequestTimeoutTotal.WithLabelValues(req.Method, req.URL.Path).Inc()
		log.Printf("⚠️ Timeout Exceeded: %s %s took %.2fs (Threshold: %.2fs)", req.Method, req.URL.Path, duration, timeoutThreshold)
	}

	// Log the request info
	log.Printf("📡 External Request: method=%s, url=%s, status=%s, duration=%.2fs",
		req.Method, req.URL.Path, statusCode, duration)

	return resp, err
}

// InitPrometheusMetrics initializes the Prometheus metrics (only once)
func InitPrometheusMetrics() {
	initOnce.Do(func() {
		log.Println("🔄 Initializing Prometheus Metrics...")

		prometheus.MustRegister(httpOutgoingRequestsTotal)
		prometheus.MustRegister(httpOutgoingRequestDuration)
		prometheus.MustRegister(httpOutgoingRequestTimeoutTotal)

		log.Println("✅ Prometheus Metrics Registered!")
	})
}
