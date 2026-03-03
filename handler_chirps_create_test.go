package main

import (
	"strings"
	"testing"
)

func TestValidateChirp_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Clean message",
			input:    "This is a nice chirp",
			expected: "This is a nice chirp",
		},
		{
			name:     "Single profanity",
			input:    "This is kerfuffle",
			expected: "This is ****",
		},
		{
			name:     "Multiple profanity",
			input:    "kerfuffle and sharbert are words",
			expected: "**** and **** are words",
		},
		{
			name:     "Case insensitive profanity",
			input:    "KERFUFFLE and Sharbert and FoRnAx",
			expected: "**** and **** and ****",
		},
		{
			name:     "Profanity at boundaries",
			input:    "fornax middle kerfuffle",
			expected: "**** middle ****",
		},
		{
			name:     "Max length valid",
			input:    "a" + strings.Repeat("b", 138) + "c", // 140 chars
			expected: "a" + strings.Repeat("b", 138) + "c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateChirp(tt.input)
			if err != nil {
				t.Fatalf("validateChirp failed: %v", err)
			}
			if result != tt.expected {
				t.Errorf("validateChirp(%q) = %q, want %q",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateChirp_TooLong(t *testing.T) {
	longChirp := strings.Repeat("a", 141) // 141 chars

	_, err := validateChirp(longChirp)
	if err == nil {
		t.Error("Expected error for chirp longer than 140 chars")
	}

	if err.Error() != "Chirp is too long" {
		t.Errorf("Expected error 'Chirp is too long', got %q", err.Error())
	}
}

func TestValidateChirp_ExactlyMaxLength(t *testing.T) {
	exactChirp := strings.Repeat("a", 140) // Exactly 140 chars

	result, err := validateChirp(exactChirp)
	if err != nil {
		t.Errorf("Should accept chirp of exactly 140 chars, got error: %v", err)
	}

	if len(result) != 140 {
		t.Errorf("Result length should be 140, got %d", len(result))
	}
}

func TestGetCleanedBody(t *testing.T) {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No bad words",
			input:    "Hello world",
			expected: "Hello world",
		},
		{
			name:     "Single bad word lowercase",
			input:    "This is kerfuffle",
			expected: "This is ****",
		},
		{
			name:     "Single bad word uppercase",
			input:    "This is KERFUFFLE",
			expected: "This is ****",
		},
		{
			name:     "Single bad word mixed case",
			input:    "This is KeRfUfFlE",
			expected: "This is ****",
		},
		{
			name:     "Multiple bad words",
			input:    "kerfuffle sharbert fornax",
			expected: "**** **** ****",
		},
		{
			name:     "Bad words mixed with good",
			input:    "I love kerfuffle but hate sharbert",
			expected: "I love **** but hate ****",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only spaces",
			input:    "   ",
			expected: "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getCleanedBody(tt.input, badWords)
			if result != tt.expected {
				t.Errorf("getCleanedBody(%q) = %q, want %q",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetCleanedBody_PreservesNonBadWords(t *testing.T) {
	badWords := map[string]struct{}{
		"bad": {},
	}

	input := "good bad better"
	expected := "good **** better"

	result := getCleanedBody(input, badWords)
	if result != expected {
		t.Errorf("Should preserve non-bad words. Got %q, want %q", result, expected)
	}
}
