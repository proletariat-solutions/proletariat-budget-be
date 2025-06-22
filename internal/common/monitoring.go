package common

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	constLabels = prometheus.Labels{"app": APPNAME}

	HTTPDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "http_request_duration_seconds",
			Help:        "Duration of HTTP requests",
			ConstLabels: constLabels,
			Buckets:     []float64{.010, .050, .100, .150, .200, .300, 0.500, 1, 5},
		}, []string{"path"})

	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "http_request_total",
			Help:        "The number of requests",
			ConstLabels: constLabels,
		}, []string{"path", "status"})

	ErrorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "rest_api_errors_total",
			Help:        "Total number of errors",
			ConstLabels: constLabels,
		}, []string{"type"})
)
