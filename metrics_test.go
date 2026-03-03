package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestHandlerReadiness(t *testing.T) {
	req := httptest.NewRequest("GET", "/admin/healthz", nil)
	w := httptest.NewRecorder()

	handlerReadiness(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if body != "OK" {
		t.Errorf("Expected body 'OK', got %q", body)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/plain; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/plain; charset=utf-8', got %q", contentType)
	}
}

func TestMetricsHandler(t *testing.T) {
	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// Sett hits til en kjent verdi
	cfg.fileserverHits.Store(42)

	req := httptest.NewRequest("GET", "/admin/metrics", nil)
	w := httptest.NewRecorder()

	cfg.metricsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "42") {
		t.Errorf("Expected body to contain '42', got %q", body)
	}

	if !strings.Contains(body, "Chirpy Admin") {
		t.Error("Expected body to contain 'Chirpy Admin'")
	}
}

func TestMiddlewareMetricsInc(t *testing.T) {
	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// Dummy handler som bare returnerer 200
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap med middleware
	wrappedHandler := cfg.middlewareMetricsInc(nextHandler)

	// Første request
	req1 := httptest.NewRequest("GET", "/app/", nil)
	w1 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w1, req1)

	if cfg.fileserverHits.Load() != 1 {
		t.Errorf("Expected hits to be 1, got %d", cfg.fileserverHits.Load())
	}

	// Andre request
	req2 := httptest.NewRequest("GET", "/app/index.html", nil)
	w2 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w2, req2)

	if cfg.fileserverHits.Load() != 2 {
		t.Errorf("Expected hits to be 2, got %d", cfg.fileserverHits.Load())
	}
}
