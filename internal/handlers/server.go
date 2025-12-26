package handlers

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"highload/internal/analytics"
	"highload/internal/cache"
	"highload/internal/metrics"
	"log"
	"net/http"
	_ "net/http/pprof"
)

type Server struct {
	redisClient *cache.RedisClient
	analyzer    *analytics.Analyzer
	handler     http.Handler
}

func NewServer(redisClient *cache.RedisClient, analyzer *analytics.Analyzer) *Server {

	srv := &Server{
		redisClient: redisClient,
		analyzer:    analyzer,
	}
	mux := http.DefaultServeMux
	mux.Handle("/prometheus", promhttp.Handler())

	mux.HandleFunc("POST /metrics", srv.metricsHandler)
	mux.HandleFunc("GET /analyze", srv.analyzeHandler)
	wrapped := metrics.LatencyMiddleware(mux)

	srv.handler = wrapped

	return srv
}

func (s *Server) Start() error {
	log.Println("Server starting on :8080")
	return http.ListenAndServe(":8080", s.handler)
}
