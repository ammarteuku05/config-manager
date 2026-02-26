package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitDB(t *testing.T) {
	// Create a temporary directory for the database file
	tempDir, err := os.MkdirTemp("", "db_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // clean up

	dbPath := filepath.Join(tempDir, "test.db")

	db, err := InitDB(dbPath)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if db == nil {
		t.Errorf("expected db to not be nil")
	}

	// Verify the file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("expected database file to be created")
	}

	// Test with invalid path (should fail)
	// NOTE: SQLite handles many invalid paths gracefully by creating them in memory
	// or current directory so it's tricky to cleanly fail without knowing the OS context perfectly.
	// An empty path usually refers to an on-disk temporary database, we can try a completely invalid path
	invalidPath := filepath.Join(tempDir, "nonexistent_dir", "test.db")
	_, errInvalid := InitDB(invalidPath)
	if errInvalid == nil {
		t.Errorf("expected error with invalid path")
	}
}
