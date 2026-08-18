package observability

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/tohjustin/aegis/internal/version"
)

const namespace = "aegis"

// Metrics holds the dedicated Prometheus registry and HTTP collectors.
type Metrics struct {
	Registry *prometheus.Registry
	Requests *prometheus.CounterVec
	Duration *prometheus.HistogramVec
}

// NewMetrics registers Go, process, HTTP, and build_info collectors on a private registry.
func NewMetrics(buckets []float64, info version.Info) (*Metrics, error) {
	reg := prometheus.NewRegistry()

	requests := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "http_requests_total",
		Help:      "Total number of HTTP requests.",
	}, []string{"method", "path", "code"})

	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "http_request_duration_seconds",
		Help:      "HTTP request duration in seconds.",
		Buckets:   buckets,
	}, []string{"method", "path"})

	buildInfo := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "build_info",
		Help:      "Build information of the running binary.",
	}, []string{"version", "git_commit", "git_tag", "git_tree_state", "go_version", "compiler", "platform"})

	collectorsToRegister := []prometheus.Collector{
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		requests,
		duration,
		buildInfo,
	}
	for _, c := range collectorsToRegister {
		if err := reg.Register(c); err != nil {
			return nil, fmt.Errorf("register collector: %w", err)
		}
	}

	buildInfo.WithLabelValues(
		info.Version,
		info.GitCommit,
		info.GitTag,
		info.GitTreeState,
		info.GoVersion,
		info.Compiler,
		info.Platform,
	).Set(1)

	return &Metrics{
		Registry: reg,
		Requests: requests,
		Duration: duration,
	}, nil
}
