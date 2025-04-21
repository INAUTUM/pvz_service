package service

import (
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	PVZCreated         prometheus.Counter
	ReceptionsCreated  prometheus.Counter
	ProductsAdded      prometheus.Counter
	HTTPRequests       *prometheus.CounterVec
	HTTPResponseTime   *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		PVZCreated: promauto.NewCounter(prometheus.CounterOpts{
			Name: "pvz_created_total",
			Help: "Total number of PVZ created",
		}),
		ReceptionsCreated: promauto.NewCounter(prometheus.CounterOpts{
			Name: "receptions_created_total",
			Help: "Total number of receptions created",
		}),
		ProductsAdded: promauto.NewCounter(prometheus.CounterOpts{
			Name: "products_added_total",
			Help: "Total number of products added",
		}),
		HTTPRequests: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		}, []string{"method", "path", "status"}),
		HTTPResponseTime: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_response_time_seconds",
			Help:    "HTTP response time distribution",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "path"}),
	}
}

func MetricsMiddleware(metrics *Metrics) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Логика сбора метрик
			return next(c)
		}
	}
}