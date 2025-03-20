package neo4jrepository

import (
	"context"
	//"fmt"
	"log"
	"os"
	"testing"

	"lk/datafoundation/crud-api/db/config"
	pb "lk/datafoundation/crud-api/lk/datafoundation/crud-api"

	"github.com/stretchr/testify/assert"
)

var repository *Neo4jRepository

// TestMain initializes the Neo4jRepository before running the tests and closes it afterward.
func TestMain(m *testing.M) {
	// Setup: Initialize the Neo4j repository with the config
	ctx := context.Background()
	cfg := &config.Neo4jConfig{
		URI:      os.Getenv("NEO4J_URI"),
		Username: os.Getenv("NEO4J_USER"),
		Password: os.Getenv("NEO4J_PASSWORD"),
	}
	log.Printf("Connecting to Neo4j at %s", cfg.URI)
	var err error
	repository, err = NewNeo4jRepository(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create Neo4j repository: %v", err)
	}

	// Run the tests
	code := m.Run()

	// Teardown: Close the repository after tests
	repository.Close(ctx)

	// Exit with the test result code
	os.Exit(code)
}

// TestCreateEntity tests the createEntity method of the Neo4jRepository
func TestCreateEntity(t *testing.T) {
	// Prepare the entity data as a map
	entity := map[string]interface{}{
		"id":          "1",
		"kind":        "Person",
		"name":        "John Doe",
		"dateCreated": "2025-03-18",
	}

	// Call the createEntity method and capture the returned entity
	createdEntity, err := repository.createEntity(context.Background(), entity)
	log.Printf("Created entity: %v", createdEntity)

	// Verify that no error occurred during creation
	assert.Nil(t, err, "Expected no error when creating an entity")

	// Verify that the returned entity has the correct values
	assert.Equal(t, "1", createdEntity["id"], "Expected entity to have the correct ID")
	assert.Equal(t, "John Doe", createdEntity["name"], "Expected entity to have the correct name")
	assert.Equal(t, "2025-03-18", createdEntity["dateCreated"], "Expected entity to have the correct creation date")

	// Optionally, if you have a 'dateEnded' field, you can check for that as well
	// assert.Nil(t, createdEntity["dateEnded"], "Expected entity to have no 'dateEnded' field")
}

// TestCreateRelationship tests the createRelationship method of the Neo4jRepository
func TestCreateRelationship(t *testing.T) {
	// Prepare the context
	ctx := context.Background()

	// Create two entities first
	entity1 := map[string]interface{}{
		"id":          "2",
		"kind":        "Person",
		"name":        "Alice",
		"dateCreated": "2025-03-18",
	}
	entity2 := map[string]interface{}{
		"id":          "3",
		"kind":        "Person",
		"name":        "Bob",
		"dateCreated": "2025-03-18",
	}

	// Create entities
	_, err := repository.createEntity(ctx, entity1)
	assert.Nil(t, err, "Expected no error when creating first entity")

	_, err = repository.createEntity(ctx, entity2)
	assert.Nil(t, err, "Expected no error when creating second entity")

	// Prepare relationship data
	relationship := &pb.Relationship{
		Id:              "101",
		RelatedEntityId: "3",
		Name:            "KNOWS",
		StartTime:       "2025-03-18",
	}

	// Create the relationship
	createdRelationship, err := repository.createRelationship(ctx, "2", relationship)
	assert.Nil(t, err, "Expected no error when creating the relationship")
	log.Printf("Created relationship: %v", createdRelationship)

	// Verify that the returned relationship has the correct values
	assert.Equal(t, "101", createdRelationship["id"], "Expected relationship to have the correct ID")
	assert.Equal(t, "2025-03-18", createdRelationship["startDate"], "Expected relationship to have the correct start date")
	assert.Equal(t, "KNOWS", createdRelationship["relationshipType"], "Expected relationship to have the correct type")

	// Optionally, check for endDate if it exists
	if relationship.EndTime != "" {
		assert.Equal(t, relationship.EndTime, createdRelationship["endDate"], "Expected relationship to have the correct end date")
	}
}

func TestReadEntity(t *testing.T) {
	// Create an entity for testing
	entity := map[string]interface{}{
		"id":          "6",
		"kind":        "Person",
		"name":        "Charlie",
		"dateCreated": "2025-03-18",
	}

	// Create the entity
	createdEntity, err := repository.createEntity(context.Background(), entity)
	assert.Nil(t, err, "Expected no error when creating the entity")
	assert.Equal(t, entity["id"], createdEntity["id"], "Expected created entity to have the correct ID")
	assert.Equal(t, entity["name"], createdEntity["name"], "Expected created entity to have the correct name")
	assert.Equal(t, entity["dateCreated"], createdEntity["dateCreated"], "Expected created entity to have the correct dateCreated")

	// Read the entity by ID
	readEntity, err := repository.ReadEntity(context.Background(), "6")
	assert.Nil(t, err, "Expected no error when reading the entity")

	// Verify the content of the entity
	assert.Equal(t, entity["id"], readEntity["id"], "Expected entity to have the correct ID")
	assert.Equal(t, entity["kind"], readEntity["kind"], "Expected entity to have the correct kind")
	assert.Equal(t, entity["name"], readEntity["name"], "Expected entity to have the correct name")
	assert.Equal(t, entity["dateCreated"], readEntity["dateCreated"], "Expected entity to have the correct dateCreated")

	// Ensure dateEnded is not present in the entity
	_, exists := readEntity["dateEnded"]
	assert.False(t, exists, "Expected entity to not have 'dateEnded' field")
}

func TestReadRelatedEntityIds(t *testing.T) {
	// Step 1: Create the entities (if not already created)
	entityMap1 := map[string]interface{}{
		"id":          "4",
		"kind":        "Person",
		"name":        "Alice",
		"dateCreated": "2025-03-18",
	}

	entityMap2 := map[string]interface{}{
		"id":          "5",
		"kind":        "Person",
		"name":        "Bob",
		"dateCreated": "2025-03-18",
	}

	// Create entities using the repository
	_, err := repository.createEntity(context.Background(), entityMap1)
	assert.Nil(t, err, "Expected no error when creating the first entity")

	_, err = repository.createEntity(context.Background(), entityMap2)
	assert.Nil(t, err, "Expected no error when creating the second entity")

	// Step 2: Create a relationship between the entities (using correct relationship creation method)
	relationship := &pb.Relationship{
		Id:              "102",        // Relationship ID
		Name:            "KNOWS",      // Relationship type (e.g., KNOWS)
		RelatedEntityId: "5",          // ID of the related entity (Bob)
		StartTime:       "2025-03-18", // Relationship start date
		EndTime:         "",           // End date (optional)
	}

	_, err = repository.createRelationship(context.Background(), "4", relationship) // "4" is the parent ID (Alice)
	assert.Nil(t, err, "Expected no error when creating the relationship")

	// Step 3: Prepare the test data for fetching related entities
	entityID := "4"             // ID of the entity to get related IDs for
	relationshipType := "KNOWS" // Relationship type
	ts := "2025-03-18"          // Timestamp (YYYY-MM-DD)

	// Step 4: Call the function to fetch related entity IDs
	relatedEntityIDs, err := repository.ReadRelatedEntityIds(context.Background(), entityID, relationshipType, ts)
	assert.Nil(t, err, "Expected no error when getting related entity IDs")
	log.Printf("Related entity IDs: %v", relatedEntityIDs)

	// Step 5: Verify the response
	assert.NotNil(t, relatedEntityIDs, "Expected related entity IDs to be returned")
	assert.Contains(t, relatedEntityIDs, "5", "Expected related entity ID 5 to be present")
}

func TestReadRelationships(t *testing.T) {
	// Create two entities
	entityMap1 := map[string]interface{}{
		"id":          "7",
		"kind":        "Person",
		"name":        "David",
		"dateCreated": "2025-03-18",
	}
	entityMap2 := map[string]interface{}{
		"id":          "8",
		"kind":        "Person",
		"name":        "Eve",
		"dateCreated": "2025-03-18",
	}

	// Create entities in the repository
	_, err := repository.createEntity(context.Background(), entityMap1)
	assert.Nil(t, err, "Expected no error when creating first entity")

	_, err = repository.createEntity(context.Background(), entityMap2)
	assert.Nil(t, err, "Expected no error when creating second entity")

	// Create a relationship between the entities
	relationship := &pb.Relationship{
		Id:              "103",
		RelatedEntityId: "8", // ID of the related entity
		Name:            "KNOWS",
		StartTime:       "2025-03-18",
	}
	_, err = repository.createRelationship(context.Background(), "7", relationship)
	assert.Nil(t, err, "Expected no error when creating the relationship")

	// Fetch relationships for entity 7
	relationships, err := repository.readRelationships(context.Background(), "7")
	assert.Nil(t, err, "Expected no error when fetching relationships")
	log.Printf("Relationships for entity 7: %v", relationships)

	// Verify that the relationship exists
	relationshipFound := false
	for _, relationship := range relationships {
		if relationship["relatedID"] == "8" {
			relationshipFound = true
			break
		}
	}

	// Assert that the relationship to the entity 8 exists
	assert.True(t, relationshipFound, "Expected relationship to include the correct related entity ID")
}

func TestGetRelationship(t *testing.T) {
	// Create two entities
	entityMap1 := map[string]interface{}{
		"id":          "9",
		"kind":        "Person",
		"name":        "David",
		"dateCreated": "2025-03-18",
	}
	entityMap2 := map[string]interface{}{
		"id":          "10",
		"kind":        "Person",
		"name":        "Eve",
		"dateCreated": "2025-03-18",
	}

	// Create entities in the repository
	_, err := repository.createEntity(context.Background(), entityMap1)
	assert.Nil(t, err, "Expected no error when creating first entity")

	_, err = repository.createEntity(context.Background(), entityMap2)
	assert.Nil(t, err, "Expected no error when creating second entity")

	// Create a relationship between the entities
	relationship := &pb.Relationship{
		Id:              "103",
		RelatedEntityId: "8", // ID of the related entity (Eve)
		Name:            "KNOWS",
		StartTime:       "2025-03-18",
	}
	_, err = repository.createRelationship(context.Background(), "7", relationship) // "7" is David's ID
	assert.Nil(t, err, "Expected no error when creating the relationship")

	// Fetch the relationship by ID
	relationshipMap, err := repository.getRelationship(context.Background(), "103")
	assert.Nil(t, err, "Expected no error when fetching the relationship")
	log.Printf("Fetched relationship: %v", relationshipMap)

	// Verify that the relationship data is correct
	assert.Equal(t, "KNOWS", relationshipMap["type"], "Expected relationship type to be KNOWS")
	assert.Equal(t, "7", relationshipMap["startEntityID"], "Expected start entity ID to be 7 (David's ID)")
	assert.Equal(t, "8", relationshipMap["endEntityID"], "Expected end entity ID to be 8 (Eve's ID)")
	assert.Equal(t, "2025-03-18", relationshipMap["startDate"], "Expected start date to be 2025-03-18")

	// Optional: Assert the endDate is nil (since it wasn't set in the creation)
	assert.Nil(t, relationshipMap["endDate"], "Expected end date to be nil")
}

func TestUpdateEntity(t *testing.T) {
	// Create a test entity
	entityData := map[string]interface{}{
		"id":          "11",
		"kind":        "Person",
		"name":        "Mary",
		"dateCreated": "2025-03-18",
	}
	_, err := repository.createEntity(context.Background(), entityData)
	assert.Nil(t, err, "Expected no error when creating entity")

	// Update the entity
	updateData := map[string]interface{}{
		"name":      "Mary Updated",
		"dateEnded": "2025-12-31",
	}

	updatedEntity, err := repository.updateEntity(context.Background(), "11", updateData)
	log.Printf("Updated entity: %v", updatedEntity)
	assert.Nil(t, err, "Expected no error when updating entity")
	assert.NotNil(t, updatedEntity, "Expected updated entity to be returned")

	// Verify that the entity was updated correctly in the return value
	assert.Equal(t, "Mary Updated", updatedEntity["name"], "Expected updated name")
	assert.Equal(t, "2025-12-31", updatedEntity["dateEnded"], "Expected updated dateEnded")

	// Fetch the entity from the database and verify
	entity, err := repository.ReadEntity(context.Background(), "11")
	log.Printf("Fetched entity: %v", entity)
	assert.Nil(t, err, "Expected no error when reading updated entity")
	assert.Equal(t, "Mary Updated", entity["name"], "Expected database to have updated name")
	assert.Equal(t, "2025-12-31", entity["dateEnded"], "Expected database to have updated dateEnded")
}

func TestUpdateRelationship(t *testing.T) {
	// Update the relationship
	updateData := map[string]interface{}{
		"endDate": "2025-12-31",
	}

	// Call the function to update the relationship
	updatedRelationship, err := repository.updateRelationship(context.Background(), "101", updateData)
	log.Printf("Updated relationship: %v", updatedRelationship)
	assert.Nil(t, err, "Expected no error when updating relationship")
	assert.NotNil(t, updatedRelationship, "Expected updated relationship to be returned")

	// Verify that the relationship was updated correctly in the return value
	assert.Equal(t, "2025-12-31", updatedRelationship["endDate"], "Expected updated endDate")

	// Fetch the relationship from the database using getRelationship
	relationship, err := repository.getRelationship(context.Background(), "101")
	log.Printf("Fetched relationship: %v", relationship)
	assert.Nil(t, err, "Expected no error when reading updated relationship")

	// Check if the relationship has the updated endDate
	assert.Equal(t, "2025-12-31", relationship["endDate"], "Expected relationship to have updated endDate")
}

func TestDeleteRelationship(t *testing.T) {
	// Ensure the relationship exists first, if needed (for test consistency)
	// You can add a create step here if the relationship doesn't exist yet

	// Call the function to delete the relationship
	err := repository.deleteRelationship(context.Background(), "101")
	assert.Nil(t, err, "Expected no error when deleting relationship")

	// Fetch the relationship to ensure it was deleted
	relationship, err := repository.getRelationship(context.Background(), "101")
	assert.NotNil(t, err, "Expected error when fetching deleted relationship")
	assert.Contains(t, err.Error(), "not found", "Expected error message to indicate relationship not found")
	assert.Nil(t, relationship, "Expected relationship to be nil after deletion")
}

func TestDeleteEntity(t *testing.T) {
	// Step 1: Test deleting an entity without relationships (ID 9)
	err := repository.deleteEntity(context.Background(), "9")
	assert.Nil(t, err, "Expected no error when deleting entity with ID 9 (no relationships)")

	// Step 2: Verify that the entity was deleted by attempting to fetch it
	_, err = repository.ReadEntity(context.Background(), "9")
	assert.NotNil(t, err, "Expected error when fetching deleted entity with ID 9")
	assert.Contains(t, err.Error(), "entity with ID 9 not found", "Expected error to indicate entity is not found")

	// Step 3: Test deleting an entity with relationships (ID 8)
	err = repository.deleteEntity(context.Background(), "8")
	assert.NotNil(t, err, "Expected error when deleting entity with ID 8 (has relationships)")

	// Step 4: Verify the error message contains information about relationships
	assert.Contains(t, err.Error(), "entity has relationships and cannot be deleted", "Expected error message to indicate relationships prevent deletion")
}
