package postgres

import (
	"os"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	// Get database URI from environment variable
	dbURI := os.Getenv("POSTGRES_TEST_DB_URI")
	if dbURI == "" {
		t.Skip("Skipping test: POSTGRES_TEST_DB_URI environment variable not set")
	}

	// Create new client using connection string
	client, err := NewClientFromDSN(dbURI)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Test basic query
	rows, err := client.DB().Query("SELECT NOW()")
	if err != nil {
		t.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Fatal("Expected at least one row")
	}

	var currentTime time.Time
	if err := rows.Scan(&currentTime); err != nil {
		t.Fatalf("Failed to scan result: %v", err)
	}

	// Verify the time is recent
	if time.Since(currentTime) > time.Minute {
		t.Errorf("Database time seems incorrect: %v", currentTime)
	}
}
