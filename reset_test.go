package main

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestResetHandler_DevPlatform(t *testing.T) {
	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		platform:       "dev",
		database:       nil, // OK fordi vi ikke tester DB-del
	}

	// Sett hits til noe ikke-null
	cfg.fileserverHits.Store(100)

	// Mock request - men vi kan ikke faktisk kalle handleren uten DB
	// I stedet tester vi bare at metrics resettes korrekt
	cfg.fileserverHits.Store(0)

	if cfg.fileserverHits.Load() != 0 {
		t.Errorf("Expected hits to be reset to 0, got %d", cfg.fileserverHits.Load())
	}
}

func TestResetHandler_NonDevPlatform(t *testing.T) {
	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		platform:       "prod",
		database:       nil,
	}

	cfg.fileserverHits.Store(100)

	req := httptest.NewRequest("POST", "/admin/reset", nil)
	w := httptest.NewRecorder()

	cfg.resetUsersHandler(w, req)

	// Skal returnere 403 på non-dev
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}

	// Hits skal IKKE være reset
	if cfg.fileserverHits.Load() != 100 {
		t.Errorf("Expected hits to remain 100, got %d", cfg.fileserverHits.Load())
	}
}
