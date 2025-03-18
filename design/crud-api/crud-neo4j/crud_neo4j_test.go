package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// setupClient initializes the Neo4jClient for testing
func setupClient(t *testing.T) *Neo4jClient {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	neo4jURI := os.Getenv("NEO4J_URI")
	neo4jUsername := os.Getenv("NEO4J_USER")
	neo4jPassword := os.Getenv("NEO4J_PASSWORD")

	client, err := NewNeo4jClient(neo4jURI, neo4jUsername, neo4jPassword)
	if err != nil {
		t.Fatalf("Failed to connect to Neo4j: %v", err)
	}

	return client
}

func TestCreateEntityDateEnded(t *testing.T) {
	client := setupClient(t)
	defer client.Close()

	entityJSON := `{
	    "id": "1234",
	    "kind": "person",
	    "name": "John Doe",
	    "dateCreated": "2025-03-17",
	    "dateEnded": "2025-03-18"
	}`

	// Step 1: Create the entity
	err := client.createEntity(entityJSON)
	if err != nil {
		t.Fatalf("Error creating entity: %v", err)
	}

	// Step 2: Retrieve the entity to verify creation
	entityID := "1234"
	retrievedJSON, err := client.readEntity(entityID)
	if err != nil {
		t.Fatalf("Error reading entity: %v", err)
	}

	// Step 3: Validate entity properties
	expected := `{"id":"1234","kind":"person","name":"John Doe","dateCreated":"2025-03-17","dateEnded":"2025-03-18"}`
	if retrievedJSON != expected {
		t.Errorf("Entity data mismatch.\nExpected: %s\nGot: %s", expected, retrievedJSON)
	}
}

func TestCreateEntity(t *testing.T) {
	client := setupClient(t)
	defer client.Close()

	entityJSON := `{
	    "id": "1235",
	    "kind": "person",
	    "name": "Jane Doe",
	    "dateCreated": "2025-03-17"
	}`

	// Step 1: Create the entity
	err := client.createEntity(entityJSON)
	if err != nil {
		t.Fatalf("Error creating entity: %v", err)
	}

	// Step 2: Retrieve the entity to verify creation
	entityID := "1235"
	retrievedJSON, err := client.readEntity(entityID)
	if err != nil {
		t.Fatalf("Error reading entity: %v", err)
	}

	// Step 3: Validate entity properties
	expected := `{"id":"1235","kind":"person","name":"Jane Doe","dateCreated":"2025-03-17"}`
	if retrievedJSON != expected {
		t.Errorf("Entity data mismatch.\nExpected: %s\nGot: %s", expected, retrievedJSON)
	}
}

func TestCreateRelationship(t *testing.T) {
	client := setupClient(t)
	defer client.Close()

	relJSON := `{
		"relationshipID": "5678",
		"parentID": "1234",
		"childID": "1235",
		"relType": "PARENT_OF",
		"startDate": "2025-03-17",
		"endDate": "2025-12-31"
	}`

	err := client.createRelationship(relJSON)
	if err != nil {
		t.Fatalf("Error creating relationship: %v", err)
	}
}

func TestReadEntity(t *testing.T) {
	client := setupClient(t)
	defer client.Close()

	entityID := "1234"
	entityJSON, err := client.readEntity(entityID)
	if err != nil {
		t.Fatalf("Error reading entity: %v", err)
	}

	fmt.Printf("Entity: %s\n", entityJSON)
}

func TestGetRelatedEntityIds(t *testing.T) {
	client := setupClient(t)
	defer client.Close()

	relatedEntityID := "1234"
	relationshipType := "PARENT_OF"
	timestamp := "2026-03-17"

	relatedIDs, err := client.getRelatedEntityIds(relatedEntityID, relationshipType, timestamp)
	if err != nil {
		t.Fatalf("Error getting related entity IDs: %v", err)
	}

	fmt.Printf("Related entity IDs: %v\n", relatedIDs)
}

func TestUpdateEntity(t *testing.T) {
	client := setupClient(t)
	defer client.Close()

	updateJSON := `{
		"id": "1234",
		"dateEnded": "2026-03-18"
	}`

	err := client.updateEntity(updateJSON)
	if err != nil {
		t.Fatalf("Error updating entity: %v", err)
	}
}

func TestUpdateRelationship(t *testing.T) {
	client := setupClient(t)
	defer client.Close()

	updateRelJSON := `{
	    "relationshipID": "5678",
	    "endDate": "2026-03-18"
	}`

	err := client.updateRelationship(updateRelJSON)
	if err != nil {
		t.Fatalf("Error updating relationship: %v", err)
	}
}

func TestDeleteRelationship(t *testing.T) {
	client := setupClient(t)
	defer client.Close()

	err := client.deleteRelationship("5678")
	if err != nil {
		t.Fatalf("Error deleting relationship: %v", err)
	}
}

func TestDeleteEntity(t *testing.T) {
	client := setupClient(t)
	defer client.Close()

	err := client.deleteEntity("1234")
	if err != nil {
		t.Fatalf("Error deleting entity: %v", err)
	}
}
