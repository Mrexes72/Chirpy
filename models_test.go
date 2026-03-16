package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUserStruct_JSONTags(t *testing.T) {
	// Test at User struct serialiserer korrekt til JSON
	user := User{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          "test@example.com",
		HashedPassword: "shouldnotappear",
		IsChirpyRed:    true,
	}

	// Serialiser til JSON
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	jsonStr := string(jsonBytes)

	// Sjekk at hashed_password IKKE er i JSON
	if strings.Contains(jsonStr, "hashed_password") {
		t.Error("JSON should not contain hashed_password field")
	}

	if strings.Contains(jsonStr, "shouldnotappear") {
		t.Error("JSON should not contain password hash value")
	}

	// Sjekk at andre felter ER der
	if !strings.Contains(jsonStr, "id") {
		t.Error("JSON should contain id field")
	}
	if !strings.Contains(jsonStr, "email") {
		t.Error("JSON should contain email field")
	}
}

func TestChirpStruct_JSONTags(t *testing.T) {
	chirp := Chirp{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      "Test chirp",
		UserID:    uuid.New(),
	}

	jsonBytes, err := json.Marshal(chirp)
	if err != nil {
		t.Fatalf("Failed to marshal chirp: %v", err)
	}

	jsonStr := string(jsonBytes)

	// Sjekk at alle felter er der
	requiredFields := []string{"id", "created_at", "updated_at", "body", "user_id"}
	for _, field := range requiredFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("JSON should contain %s field", field)
		}
	}
}
