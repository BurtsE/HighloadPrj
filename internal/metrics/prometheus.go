package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of requests",
		},
		[]string{"handler"},
	)
	AnomaliesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "anomalies_total",
			Help: "Total number of anomalies detected",
		},
	)
	LatencyHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "request_duration_seconds",
			Help: "Request duration in seconds",
		},
	)
)

func init() {
	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(AnomaliesTotal)
	prometheus.MustRegister(LatencyHistogram)
}

func LatencyMiddleware(next http.Handler) http.Handler { // Для Gorilla Mux (http.Handler)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Таймер для latency
		defer func() {
			LatencyHistogram.Observe(time.Since(start).Seconds())
		}()
		next.ServeHTTP(w, r)
	})
}
