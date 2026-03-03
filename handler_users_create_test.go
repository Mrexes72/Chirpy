package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestCreateUserHandler_Success(t *testing.T) {
	// TODO: Dette krever mock database
	// For nå, skip denne
	t.Skip("Requires database mock")
}

func TestRespondWithJSON(t *testing.T) {
	// Test helper-funksjonen
	w := httptest.NewRecorder()

	type testData struct {
		Message string `json:"message"`
		Count   int    `json:"count"`
	}

	data := testData{
		Message: "test",
		Count:   42,
	}

	respondWithJSON(w, 200, data)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var result testData
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result.Message != "test" || result.Count != 42 {
		t.Errorf("Response data mismatch: %+v", result)
	}
}

func TestRespondWithError(t *testing.T) {
	w := httptest.NewRecorder()

	respondWithError(w, 404, "Not found", nil)

	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	var result map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result["error"] != "Not found" {
		t.Errorf("Expected error message 'Not found', got %s", result["error"])
	}
}
