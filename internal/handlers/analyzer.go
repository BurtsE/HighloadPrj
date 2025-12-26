package handlers

import (
	"encoding/json"
	"highload/internal/metrics"
	"net/http"
)

func (s *Server) analyzeHandler(w http.ResponseWriter, r *http.Request) {
	go metrics.RequestsTotal.WithLabelValues("analyze").Inc()
	avg, anomalies := s.analyzer.GetStats()
	response := map[string]interface{}{
		"rolling_avg_rps":  avg,
		"recent_anomalies": anomalies,
	}
	_ = json.NewEncoder(w).Encode(response)

}
