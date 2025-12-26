package handlers

import (
	"encoding/json"
	"highload/internal/cache"
	"highload/internal/metrics"
	"highload/internal/model"
	"log"
	"net/http"
)

func (s *Server) metricsHandler(w http.ResponseWriter, r *http.Request) {
	go metrics.RequestsTotal.WithLabelValues("metrics").Inc()

	var m model.Metric
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Сохраняем метрику в Redis
	go func() {
		err := cache.SaveMetric(s.redisClient, m)
		if err != nil {
			log.Printf("Failed to save metric: %v", err)
		}
	}()

	// Обновляем аналитику
	go s.analyzer.Update(m)

	w.WriteHeader(http.StatusAccepted)
}
