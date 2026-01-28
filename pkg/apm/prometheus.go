package apm

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus provides lightweight application performance monitoring via
// Prometheus metrics for HTTP request/response cycles.
type Prometheus struct {
	service string
	reg     *prometheus.Registry

	inFlight prometheus.Gauge

	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewPrometheus(service string) *Prometheus {
	if service == "" {
		service = "go-url-shortener"
	}

	reg := prometheus.NewRegistry()
	_ = reg.Register(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	_ = reg.Register(prometheus.NewGoCollector())

	p := &Prometheus{
		service: service,
		reg:     reg,
		inFlight: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "apm",
			Subsystem: "http",
			Name:      "in_flight_requests",
			Help:      "Current number of in-flight HTTP requests.",
			ConstLabels: prometheus.Labels{
				"service": service,
			},
		}),
		requestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "apm",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests processed.",
			ConstLabels: prometheus.Labels{
				"service": service,
			},
		}, []string{"method", "route", "status"}),
		requestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "apm",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "HTTP request duration in seconds.",
			Buckets:   prometheus.DefBuckets,
			ConstLabels: prometheus.Labels{
				"service": service,
			},
		}, []string{"method", "route", "status"}),
	}

	reg.MustRegister(p.inFlight, p.requestsTotal, p.requestDuration)
	return p
}

func (p *Prometheus) Registry() *prometheus.Registry {
	return p.reg
}

func (p *Prometheus) Handler() http.Handler {
	return promhttp.HandlerFor(p.reg, promhttp.HandlerOpts{})
}

// GinMiddleware instruments requests with Prometheus metrics.
func (p *Prometheus) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		p.inFlight.Inc()
		defer p.inFlight.Dec()

		c.Next()

		route := c.FullPath()
		if route == "" {
			// FullPath is empty when no route matched.
			route = "unmatched"
		}

		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		p.requestsTotal.WithLabelValues(method, route, status).Inc()
		p.requestDuration.WithLabelValues(method, route, status).Observe(time.Since(start).Seconds())
	}
}
