package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Metric struct {
	ID        int       `json:"id"`
	EventName string    `json:"event_name"`
	Count     int       `json:"count"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MetricStore struct {
	mu      sync.RWMutex
	metrics map[string]*Metric
	nextID  int
}

func NewMetricStore() *MetricStore {
	return &MetricStore{metrics: make(map[string]*Metric), nextID: 1}
}

func (s *MetricStore) Increment(eventName string) *Metric {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.metrics[eventName]
	if !ok {
		m = &Metric{ID: s.nextID, EventName: eventName}
		s.nextID++
		s.metrics[eventName] = m
	}
	m.Count++
	m.UpdatedAt = time.Now().UTC()
	return m
}

func (s *MetricStore) List() []*Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Metric, 0, len(s.metrics))
	for _, m := range s.metrics {
		result = append(result, m)
	}
	return result
}

var store = NewMetricStore()

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":    "ok",
		"service":   "analytics",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

type trackRequest struct {
	EventName string `json:"event_name"`
}

func trackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var req trackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
		return
	}
	if req.EventName == "" {
		http.Error(w, `{"error":"event_name is required"}`, http.StatusBadRequest)
		return
	}
	m := store.Increment(req.EventName)
	log.Printf("Tracked event: %s count=%d", req.EventName, m.Count)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := store.List()
	log.Printf("Listing %d metrics", len(metrics))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"metrics": metrics,
		"total":   len(metrics),
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/track", trackHandler)
	mux.HandleFunc("/metrics", metricsHandler)

	log.Printf("Analytics service starting on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
