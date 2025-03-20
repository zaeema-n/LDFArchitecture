package neo4jrepository

import (
	"context"
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

// TestCreateEntity tests the CreateGraphEntity method of the Neo4jRepository
func TestCreateEntity(t *testing.T) {
	// Prepare the entity data as a map
	entity := map[string]interface{}{
		"Id":         "1",
		"Kind":       "Person",
		"Name":       "John Doe",
		"Created":    "2025-03-18",
		"Terminated": nil,
	}

	// Call the CreateGraphEntity method and capture the returned entity
	createdEntity, err := repository.CreateGraphEntity(context.Background(), entity)
	log.Printf("Created entity: %v", createdEntity)

	// Verify that no error occurred during creation
	assert.Nil(t, err, "Expected no error when creating an entity")

	// Verify that the returned entity has the correct values
	assert.Equal(t, "1", createdEntity["Id"], "Expected entity to have the correct Id")
	assert.Equal(t, "John Doe", createdEntity["Name"], "Expected entity to have the correct Name")
	assert.Equal(t, "2025-03-18", createdEntity["Created"], "Expected entity to have the correct Created date")
	assert.Nil(t, createdEntity["Terminated"], "Expected entity to have no Terminated field")
}

// TestCreateRelationship tests the CreateRelationship method of the Neo4jRepository
func TestCreateRelationship(t *testing.T) {
	// Prepare the context
	ctx := context.Background()

	// Create two entities first
	entity1 := map[string]interface{}{
		"Id":      "2",
		"Kind":    "Person",
		"Name":    "Alice",
		"Created": "2025-03-18",
	}
	entity2 := map[string]interface{}{
		"Id":      "3",
		"Kind":    "Person",
		"Name":    "Bob",
		"Created": "2025-03-18",
	}

	// Create entities
	_, err := repository.CreateGraphEntity(ctx, entity1)
	assert.Nil(t, err, "Expected no error when creating first entity")

	_, err = repository.CreateGraphEntity(ctx, entity2)
	assert.Nil(t, err, "Expected no error when creating second entity")

	// Prepare relationship data
	relationship := &pb.Relationship{
		Id:              "101",
		RelatedEntityId: "3",
		Name:            "KNOWS",
		StartTime:       "2025-03-18",
	}

	// Create the relationship
	createdRelationship, err := repository.CreateRelationship(ctx, "2", relationship)
	assert.Nil(t, err, "Expected no error when creating the relationship")
	log.Printf("Created relationship: %v", createdRelationship)

	// Verify that the returned relationship has the correct values
	assert.Equal(t, "101", createdRelationship["Id"], "Expected relationship to have the correct Id")
	assert.Equal(t, "2025-03-18", createdRelationship["Created"], "Expected relationship to have the correct Created date")
	assert.Equal(t, "KNOWS", createdRelationship["relationshipType"], "Expected relationship to have the correct type")
}

// TestReadEntity tests the ReadGraphEntity method of the Neo4jRepository
func TestReadEntity(t *testing.T) {
	// Create an entity for testing
	entity := map[string]interface{}{
		"Id":      "6",
		"Kind":    "Person",
		"Name":    "Charlie",
		"Created": "2025-03-18",
	}

	// Create the entity
	createdEntity, err := repository.CreateGraphEntity(context.Background(), entity)
	assert.Nil(t, err, "Expected no error when creating the entity")
	assert.Equal(t, entity["Id"], createdEntity["Id"], "Expected created entity to have the correct Id")
	assert.Equal(t, entity["Name"], createdEntity["Name"], "Expected created entity to have the correct Name")
	assert.Equal(t, entity["Created"], createdEntity["Created"], "Expected created entity to have the correct Created date")

	// Read the entity by Id
	readEntity, err := repository.ReadGraphEntity(context.Background(), "6")
	assert.Nil(t, err, "Expected no error when reading the entity")

	// Verify the content of the entity
	assert.Equal(t, entity["Id"], readEntity["Id"], "Expected entity to have the correct Id")
	assert.Equal(t, entity["Kind"], readEntity["Kind"], "Expected entity to have the correct Kind")
	assert.Equal(t, entity["Name"], readEntity["Name"], "Expected entity to have the correct Name")
	assert.Equal(t, entity["Created"], readEntity["Created"], "Expected entity to have the correct Created date")
}

// TestReadRelatedEntityIds tests the ReadRelatedGraphEntityIds method of the Neo4jRepository
func TestReadRelatedEntityIds(t *testing.T) {
	// Create two entities
	entity1 := map[string]interface{}{
		"Id":      "4",
		"Kind":    "Person",
		"Name":    "Alice",
		"Created": "2025-03-18",
	}
	entity2 := map[string]interface{}{
		"Id":      "5",
		"Kind":    "Person",
		"Name":    "Bob",
		"Created": "2025-03-18",
	}

	// Create entities
	_, err := repository.CreateGraphEntity(context.Background(), entity1)
	assert.Nil(t, err, "Expected no error when creating the first entity")

	_, err = repository.CreateGraphEntity(context.Background(), entity2)
	assert.Nil(t, err, "Expected no error when creating the second entity")

	// Create a relationship between the entities
	relationship := &pb.Relationship{
		Id:              "102",
		Name:            "KNOWS",
		RelatedEntityId: "5",
		StartTime:       "2025-03-18",
	}

	_, err = repository.CreateRelationship(context.Background(), "4", relationship)
	assert.Nil(t, err, "Expected no error when creating the relationship")

	// Step 3: Prepare the test data for fetching related entities
	entityID := "4"             // ID of the entity to get related IDs for
	relationshipType := "KNOWS" // Relationship type
	ts := "2025-03-18"          // Timestamp (YYYY-MM-DD)

	// Step 4: Call the function to fetch related entity IDs
	relatedEntityIDs, err := repository.ReadRelatedGraphEntityIds(context.Background(), entityID, relationshipType, ts)
	assert.Nil(t, err, "Expected no error when getting related entity IDs")
	log.Printf("Related entity IDs: %v", relatedEntityIDs)

	// Step 5: Verify the response
	assert.NotNil(t, relatedEntityIDs, "Expected related entity IDs to be returned")
	assert.Contains(t, relatedEntityIDs, "5", "Expected related entity ID 5 to be present")
}

func TestReadRelationships(t *testing.T) {
	// Create two entities
	entityMap1 := map[string]interface{}{
		"Id":      "7",
		"Kind":    "Person",
		"Name":    "David",
		"Created": "2025-03-18",
	}
	entityMap2 := map[string]interface{}{
		"Id":      "8",
		"Kind":    "Person",
		"Name":    "Eve",
		"Created": "2025-03-18",
	}

	// Create entities in the repository
	_, err := repository.CreateGraphEntity(context.Background(), entityMap1)
	assert.Nil(t, err, "Expected no error when creating first entity")

	_, err = repository.CreateGraphEntity(context.Background(), entityMap2)
	assert.Nil(t, err, "Expected no error when creating second entity")

	// Create a relationship between the entities
	relationship := &pb.Relationship{
		Id:              "103",
		RelatedEntityId: "8", // ID of the related entity
		Name:            "KNOWS",
		StartTime:       "2025-03-18",
	}
	_, err = repository.CreateRelationship(context.Background(), "7", relationship)
	assert.Nil(t, err, "Expected no error when creating the relationship")

	// Fetch relationships for entity 7
	relationships, err := repository.ReadRelationships(context.Background(), "7")
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

func TestReadRelationship(t *testing.T) {
	// Create two entities
	entityMap1 := map[string]interface{}{
		"Id":      "9",
		"Kind":    "Person",
		"Name":    "David",
		"Created": "2025-03-18",
	}
	entityMap2 := map[string]interface{}{
		"Id":      "10",
		"Kind":    "Person",
		"Name":    "Eve",
		"Created": "2025-03-18",
	}

	// Create entities in the repository
	_, err := repository.CreateGraphEntity(context.Background(), entityMap1)
	assert.Nil(t, err, "Expected no error when creating first entity")

	_, err = repository.CreateGraphEntity(context.Background(), entityMap2)
	assert.Nil(t, err, "Expected no error when creating second entity")

	// Create a relationship between the entities
	relationship := &pb.Relationship{
		Id:              "103",
		RelatedEntityId: "8", // ID of the related entity (Eve)
		Name:            "KNOWS",
		StartTime:       "2025-03-18",
	}
	_, err = repository.CreateRelationship(context.Background(), "7", relationship) // "7" is David's ID
	assert.Nil(t, err, "Expected no error when creating the relationship")

	// Fetch the relationship by ID
	relationshipMap, err := repository.ReadRelationship(context.Background(), "103")
	assert.Nil(t, err, "Expected no error when fetching the relationship")
	log.Printf("Fetched relationship: %v", relationshipMap)

	// Verify that the relationship data is correct
	assert.Equal(t, "KNOWS", relationshipMap["type"], "Expected relationship type to be KNOWS")
	assert.Equal(t, "7", relationshipMap["startEntityID"], "Expected start entity ID to be 7 (David's ID)")
	assert.Equal(t, "8", relationshipMap["endEntityID"], "Expected end entity ID to be 8 (Eve's ID)")
	assert.Equal(t, "2025-03-18", relationshipMap["Created"], "Expected start date to be 2025-03-18")

	// Optional: Assert the endDate is nil (since it wasn't set in the creation)
	assert.Nil(t, relationshipMap["Terminated"], "Expected end date to be nil")
}

func TestUpdateEntity(t *testing.T) {
	// Create a test entity
	entityData := map[string]interface{}{
		"Id":      "11",
		"Kind":    "Person",
		"Name":    "Mary",
		"Created": "2025-03-18",
	}
	_, err := repository.CreateGraphEntity(context.Background(), entityData)
	assert.Nil(t, err, "Expected no error when creating entity")

	// Update the entity
	updateData := map[string]interface{}{
		"Name":       "Mary Updated",
		"Terminated": "2025-12-31",
	}

	updatedEntity, err := repository.UpdateGraphEntity(context.Background(), "11", updateData)
	log.Printf("Updated entity: %v", updatedEntity)
	assert.Nil(t, err, "Expected no error when updating entity")
	assert.NotNil(t, updatedEntity, "Expected updated entity to be returned")

	// Verify that the entity was updated correctly in the return value
	assert.Equal(t, "Mary Updated", updatedEntity["Name"], "Expected updated name")
	assert.Equal(t, "2025-12-31", updatedEntity["Terminated"], "Expected updated dateEnded")

	// Fetch the entity from the database and verify
	entity, err := repository.ReadGraphEntity(context.Background(), "11")
	log.Printf("Fetched entity: %v", entity)
	assert.Nil(t, err, "Expected no error when reading updated entity")
	assert.Equal(t, "Mary Updated", entity["Name"], "Expected database to have updated name")
	assert.Equal(t, "2025-12-31", entity["Terminated"], "Expected database to have updated dateEnded")
}

func TestUpdateRelationship(t *testing.T) {
	// Update the relationship
	updateData := map[string]interface{}{
		"Terminated": "2025-12-31",
	}

	// Call the function to update the relationship
	updatedRelationship, err := repository.UpdateRelationship(context.Background(), "101", updateData)
	log.Printf("Updated relationship: %v", updatedRelationship)
	assert.Nil(t, err, "Expected no error when updating relationship")
	assert.NotNil(t, updatedRelationship, "Expected updated relationship to be returned")

	// Verify that the relationship was updated correctly in the return value
	assert.Equal(t, "2025-12-31", updatedRelationship["Terminated"], "Expected updated endDate")

	// Fetch the relationship from the database using getRelationship
	relationship, err := repository.ReadRelationship(context.Background(), "101")
	log.Printf("Fetched relationship: %v", relationship)
	assert.Nil(t, err, "Expected no error when reading updated relationship")

	// Check if the relationship has the updated endDate
	assert.Equal(t, "2025-12-31", relationship["Terminated"], "Expected relationship to have updated endDate")
}

func TestDeleteRelationship(t *testing.T) {
	// Ensure the relationship exists first, if needed (for test consistency)
	// You can add a create step here if the relationship doesn't exist yet

	// Call the function to delete the relationship
	err := repository.DeleteRelationship(context.Background(), "101")
	assert.Nil(t, err, "Expected no error when deleting relationship")

	// Fetch the relationship to ensure it was deleted
	relationship, err := repository.ReadRelationship(context.Background(), "101")
	assert.NotNil(t, err, "Expected error when fetching deleted relationship")
	assert.Contains(t, err.Error(), "not found", "Expected error message to indicate relationship not found")
	assert.Nil(t, relationship, "Expected relationship to be nil after deletion")
}

func TestDeleteEntity(t *testing.T) {
	// Create a test entity
	entity := map[string]interface{}{
		"Id":      "12",
		"Kind":    "Person",
		"Name":    "John Smith",
		"Created": "2025-03-18",
	}
	_, err := repository.CreateGraphEntity(context.Background(), entity)
	assert.Nil(t, err, "Expected no error when creating entity")

	// Step 2: Verify that the entity was deleted by attempting to fetch it
	// _, err = repository.ReadGraphEntity(context.Background(), "9")
	// assert.NotNil(t, err, "Expected error when fetching deleted entity with ID 9")
	// assert.Contains(t, err.Error(), "entity with ID 9 not found", "Expected error to indicate entity is not found")

	err = repository.DeleteGraphEntity(context.Background(), "12")
	assert.Nil(t, err, "Expected no error when deleting entity")

	// Verify the entity was deleted
	_, err = repository.ReadGraphEntity(context.Background(), "12")
	assert.NotNil(t, err, "Expected error when fetching deleted entity")
	assert.Contains(t, err.Error(), "not found", "Expected error message to indicate entity not found")

	// Step 3: Test deleting an entity with relationships (ID 8)
	err = repository.DeleteGraphEntity(context.Background(), "8")
	assert.NotNil(t, err, "Expected error when deleting entity with ID 8 (has relationships)")

	// Step 4: Verify the error message contains information about relationships
	assert.Contains(t, err.Error(), "entity has relationships and cannot be deleted", "Expected error message to indicate relationships prevent deletion")
}
