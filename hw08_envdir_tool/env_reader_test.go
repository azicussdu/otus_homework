package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	// Temporary directory for testing.
	tempDir := t.TempDir()

	createTestFile(t, tempDir, "FOO", "value\n")
	createTestFile(t, tempDir, "BOO", " another value  \n")
	createTestFile(t, tempDir, "EMPTY", " ")
	createTestFile(t, tempDir, "UNSET", "")
	createTestFile(t, tempDir, "FOO2", "first line\nsecond line\n")

	// Run ReadDir function on the test directory.
	env, err := ReadDir(tempDir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}

	// Define expected results.
	expected := Environment{
		"FOO":   {Value: "value", NeedRemove: false},
		"BOO":   {Value: " another value", NeedRemove: false},
		"EMPTY": {Value: "", NeedRemove: false},
		"UNSET": {Value: "", NeedRemove: true},
		"FOO2":  {Value: "first line", NeedRemove: false},
	}

	// Compare actual and expected results.
	for key, expectedValue := range expected {
		gotValue, exists := env[key]
		if !exists {
			t.Errorf("expected key %s not found in result", key)
		}
		if gotValue != expectedValue {
			t.Errorf("for key %s, expected %v, got %v", key, expectedValue, gotValue)
		}
	}
}

func createTestFile(t *testing.T, dir, name, content string) {
	t.Helper()
	filePath := filepath.Join(dir, name)
	// it shows me gofumpt linter error, but everything is correct here, so I disabled it.
	//nolint:gofumpt
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file %s: %v", name, err)
	}
}
