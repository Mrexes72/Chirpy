package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup før alle tester kjører

	// Kjør testene
	code := m.Run()

	// Cleanup etter tester

	os.Exit(code)
}
