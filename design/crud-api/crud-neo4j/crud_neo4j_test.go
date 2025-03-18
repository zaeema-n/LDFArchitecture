package main

import (
	"os"
	"testing"

	"encoding/json"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var client *Neo4jClient

func TestMain(m *testing.M) {
	// Load environment variables from the .env file
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	// Connect to the same Neo4j database
	neo4jURI := os.Getenv("NEO4J_URI")
	neo4jUsername := os.Getenv("NEO4J_USER")
	neo4jPassword := os.Getenv("NEO4J_PASSWORD")

	// Initialize Neo4j client
	var err2 error
	client, err2 = NewNeo4jClient(neo4jURI, neo4jUsername, neo4jPassword)
	if err2 != nil {
		panic("Failed to connect to Neo4j")
	}
	defer client.Close()

	// Run tests
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCreateEntity(t *testing.T) {
	entityJSON := `{
		"id": "1",
		"kind": "Person",
		"name": "John Doe",
		"dateCreated": "2025-03-18"
	}`

	err := client.createEntity(entityJSON)
	assert.Nil(t, err, "Expected no error when creating an entity")

	// Verify that the entity was created
	entity, err := client.readEntity("1")
	assert.Nil(t, err, "Expected no error when reading the created entity")
	assert.Contains(t, entity, `"name":"John Doe"`, "Expected entity to have the correct name")
}

func TestReadEntity(t *testing.T) {
	// Create an entity for testing
	entityJSON := `{
		"id": "6",
		"kind": "Person",
		"name": "Charlie",
		"dateCreated": "2025-03-18"
	}`

	err := client.createEntity(entityJSON)
	assert.Nil(t, err, "Expected no error when creating the entity")

	// Fetch the entity by ID
	entity, err := client.readEntity("6")
	assert.Nil(t, err, "Expected no error when reading the entity")

	// Verify the content of the entity
	assert.Contains(t, entity, `"name":"Charlie"`, "Expected entity to have the correct name")
	assert.Contains(t, entity, `"dateCreated":"2025-03-18"`, "Expected entity to have the correct dateCreated")
}

func TestCreateRelationship(t *testing.T) {
	// Create two entities first
	entityJSON1 := `{
		"id": "2",
		"kind": "Person",
		"name": "Alice",
		"dateCreated": "2025-03-18"
	}`
	entityJSON2 := `{
		"id": "3",
		"kind": "Person",
		"name": "Bob",
		"dateCreated": "2025-03-18"
	}`

	err := client.createEntity(entityJSON1)
	assert.Nil(t, err, "Expected no error when creating first entity")

	err = client.createEntity(entityJSON2)
	assert.Nil(t, err, "Expected no error when creating second entity")

	// Create a relationship between the entities
	relationshipJSON := `{
		"relationshipID": "101",
		"parentID": "2",
		"childID": "3",
		"relType": "KNOWS",
		"startDate": "2025-03-18"
	}`

	err = client.createRelationship(relationshipJSON)
	assert.Nil(t, err, "Expected no error when creating a relationship")

	// Verify relationship creation
	rawRelationships, err := client.getRelationships("2")
	assert.Nil(t, err, "Expected no error when fetching relationships")

	// Parse the raw JSON response into a slice of maps
	var relationships []map[string]interface{}
	err = json.Unmarshal([]byte(rawRelationships), &relationships)
	assert.Nil(t, err, "Expected no error when unmarshalling relationships")

	// Check that one of the relationships contains the expected relatedID
	found := false
	for _, relationship := range relationships {
		if relationship["relatedID"] == "3" {
			found = true
			break
		}
	}

	assert.True(t, found, "Expected relationship to include the correct related entity ID")
}

func TestGetRelationships(t *testing.T) {
	// Create two entities
	entityJSON1 := `{
		"id": "7",
		"kind": "Person",
		"name": "David",
		"dateCreated": "2025-03-18"
	}`
	entityJSON2 := `{
		"id": "8",
		"kind": "Person",
		"name": "Eve",
		"dateCreated": "2025-03-18"
	}`

	err := client.createEntity(entityJSON1)
	assert.Nil(t, err, "Expected no error when creating first entity")

	err = client.createEntity(entityJSON2)
	assert.Nil(t, err, "Expected no error when creating second entity")

	// Create a relationship between the entities
	relationshipJSON := `{
		"relationshipID": "103",
		"parentID": "7",
		"childID": "8",
		"relType": "KNOWS",
		"startDate": "2025-03-18"
	}`

	err = client.createRelationship(relationshipJSON)
	assert.Nil(t, err, "Expected no error when creating the relationship")

	// Fetch the relationships for entity 7
	rawRelationships, err := client.getRelationships("7")
	assert.Nil(t, err, "Expected no error when fetching relationships")

	// Parse the raw JSON response into a slice of maps
	var relationships []map[string]interface{}
	err = json.Unmarshal([]byte(rawRelationships), &relationships)
	assert.Nil(t, err, "Expected no error when unmarshalling relationships")

	// Verify that the relationship exists
	relationshipFound := false
	for _, relationship := range relationships {
		if relationship["relatedID"] == "8" {
			relationshipFound = true
			break
		}
	}

	assert.True(t, relationshipFound, "Expected relationship to include the correct related entity ID")
}

func TestGetRelatedEntityIds(t *testing.T) {
	// Step 1: Create the entities (if not already created)
	entityJSON1 := `{
		"id": "4",
		"kind": "Person",
		"name": "Alice",
		"dateCreated": "2025-03-18"
	}`
	entityJSON2 := `{
		"id": "5",
		"kind": "Person",
		"name": "Bob",
		"dateCreated": "2025-03-18"
	}`

	err := client.createEntity(entityJSON1)
	assert.Nil(t, err, "Expected no error when creating the first entity")

	err = client.createEntity(entityJSON2)
	assert.Nil(t, err, "Expected no error when creating the second entity")

	// Step 2: Create a relationship between the entities (if not already created)
	relationshipJSON := `{
		"relationshipID": "102",
		"parentID": "4",
		"childID": "5",
		"relType": "KNOWS",
		"startDate": "2025-03-18"
	}`

	err = client.createRelationship(relationshipJSON)
	assert.Nil(t, err, "Expected no error when creating the relationship")

	// Step 3: Prepare the test data for fetching related entities
	entityID := "4"         // ID of the entity to get related IDs for
	relationship := "KNOWS" // Relationship type
	ts := "2025-03-18"      // Timestamp (YYYY-MM-DD)

	// Step 4: Call the function to fetch related entity IDs
	relatedEntityIDs, err := client.getRelatedEntityIds(entityID, relationship, ts)
	assert.Nil(t, err, "Expected no error when getting related entity IDs")

	// Step 5: Verify the response
	assert.NotNil(t, relatedEntityIDs, "Expected related entity IDs to be returned")
	assert.Contains(t, relatedEntityIDs, "5", "Expected related entity ID 5 to be present")
}

func TestUpdateRelationship(t *testing.T) {
	// Prepare the update JSON
	updateJSON := `{
		"relationshipID": "102",
		"endDate": "2025-12-31"
	}`

	// Call the function to update the relationship
	err := client.updateRelationship(updateJSON)
	assert.Nil(t, err, "Expected no error when updating the relationship")

	// Verify the relationship has been updated (check if the endDate exists)
	rawRelationships, err := client.getRelationships("5")
	assert.Nil(t, err, "Expected no error when fetching relationships")

	// Parse the raw JSON response into a slice of maps
	var relationships []map[string]interface{}
	err = json.Unmarshal([]byte(rawRelationships), &relationships)
	assert.Nil(t, err, "Expected no error when unmarshalling relationships")

	// Check if the updated relationship contains the new endDate
	relationshipUpdated := false
	for _, relationship := range relationships {
		if relationship["relationshipID"] == "102" {
			endDate, exists := relationship["endDate"]
			if exists && endDate == "2025-12-31" {
				relationshipUpdated = true
			}
		}
	}

	assert.True(t, relationshipUpdated, "Expected relationship to have updated endDate")
}

func TestUpdateEntity(t *testing.T) {
	updateJSON := `{
		"id": "1",
		"name": "John Updated",
		"dateEnded": "2025-12-31"
	}`

	err := client.updateEntity(updateJSON)
	assert.Nil(t, err, "Expected no error when updating entity")

	// Verify that the entity was updated
	entity, err := client.readEntity("1")
	assert.Nil(t, err, "Expected no error when reading updated entity")
	assert.Contains(t, entity, `"name":"John Updated"`, "Expected entity to have updated name")
	assert.Contains(t, entity, `"dateEnded":"2025-12-31"`, "Expected entity to have updated dateEnded")
}

func TestDeleteEntity(t *testing.T) {
	// Ensure no relationships exist for the entity
	err := client.deleteRelationship("101")
	assert.Nil(t, err, "Expected no error when deleting the relationship")

	// Now delete the entity
	err = client.deleteEntity("1")
	assert.Nil(t, err, "Expected no error when deleting entity")

	// Verify that the entity is deleted
	entity, err := client.readEntity("1")
	assert.NotNil(t, err, "Expected error when reading deleted entity")
	assert.Empty(t, entity, "Expected entity to be deleted")
}

func TestDeleteRelationship(t *testing.T) {
	// Ensure the relationship exists first
	relationshipJSON := `{
		"relationshipID": "101",
		"parentID": "2",
		"childID": "3",
		"relType": "KNOWS",
		"startDate": "2025-03-18"
	}`

	// Fetch current relationships of the parent
	rawRelationships, err := client.getRelationships("2")
	if err != nil {
		t.Fatalf("Expected no error when fetching relationships: %v", err)
	}

	// Parse the raw JSON response into a slice of maps
	var relationships []map[string]interface{}
	err = json.Unmarshal([]byte(rawRelationships), &relationships)
	if err != nil {
		t.Fatalf("Expected no error when unmarshalling relationships: %v", err)
	}

	// Check if relationship with ID "101" already exists
	relationshipExists := false
	for _, relationship := range relationships {
		if relationship["relationshipID"] == "101" {
			relationshipExists = true
			break
		}
	}

	// If the relationship does not exist, create it
	if !relationshipExists {
		err := client.createRelationship(relationshipJSON)
		assert.Nil(t, err, "Expected no error when creating the relationship")
	}

	// Now delete the relationship
	err = client.deleteRelationship("101")
	if err != nil && err.Error() != "relationship with ID 101 does not exist" {
		t.Fatalf("Expected no error when deleting the relationship, got: %v", err)
	}

	// Verify the relationship is deleted
	rawRelationships, err = client.getRelationships("2")
	assert.Nil(t, err, "Expected no error when fetching relationships")

	// Parse the raw JSON response into a slice of maps
	err = json.Unmarshal([]byte(rawRelationships), &relationships)
	assert.Nil(t, err, "Expected no error when unmarshalling relationships")

	assert.NotContains(t, rawRelationships, `"relationshipID":"101"`, "Expected relationship to be deleted")
}
