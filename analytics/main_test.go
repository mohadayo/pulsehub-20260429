package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	healthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["status"] != "ok" {
		t.Fatalf("expected ok, got %s", resp["status"])
	}
	if resp["service"] != "analytics" {
		t.Fatalf("expected analytics, got %s", resp["service"])
	}
}

func TestTrackHandler(t *testing.T) {
	store = NewMetricStore()
	body, _ := json.Marshal(trackRequest{EventName: "click"})
	req := httptest.NewRequest(http.MethodPost, "/track", bytes.NewReader(body))
	w := httptest.NewRecorder()
	trackHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var m Metric
	json.NewDecoder(w.Body).Decode(&m)
	if m.EventName != "click" {
		t.Fatalf("expected click, got %s", m.EventName)
	}
	if m.Count != 1 {
		t.Fatalf("expected count 1, got %d", m.Count)
	}
}

func TestTrackHandlerIncrement(t *testing.T) {
	store = NewMetricStore()
	body, _ := json.Marshal(trackRequest{EventName: "view"})
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/track", bytes.NewReader(body))
		w := httptest.NewRecorder()
		trackHandler(w, req)
	}
	metrics := store.List()
	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}
	if metrics[0].Count != 3 {
		t.Fatalf("expected count 3, got %d", metrics[0].Count)
	}
}

func TestTrackHandlerEmptyName(t *testing.T) {
	body, _ := json.Marshal(trackRequest{EventName: ""})
	req := httptest.NewRequest(http.MethodPost, "/track", bytes.NewReader(body))
	w := httptest.NewRecorder()
	trackHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestTrackHandlerWrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/track", nil)
	w := httptest.NewRecorder()
	trackHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestMetricsHandler(t *testing.T) {
	store = NewMetricStore()
	store.Increment("a")
	store.Increment("b")

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	metricsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if int(resp["total"].(float64)) != 2 {
		t.Fatalf("expected 2 metrics, got %v", resp["total"])
	}
}
