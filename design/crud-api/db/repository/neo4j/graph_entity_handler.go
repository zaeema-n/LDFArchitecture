package neo4jrepository

import (
	"context"
	"fmt"
	"log"

	pb "lk/datafoundation/crud-api/lk/datafoundation/crud-api" // Replace with your actual protobuf package

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// GetEntityDetailsFromNeo4j retrieves entity information from Neo4j database
func (repo *Neo4jRepository) GetGraphEntity(ctx context.Context, entityId string) (*pb.Kind, *pb.TimeBasedValue, string, string, error) {
	// Try to get additional entity information from Neo4j
	var kind *pb.Kind
	var name *pb.TimeBasedValue
	var created string
	var terminated string

	// Attempt to read from Neo4j, but don't fail if it's not available
	entityMap, err := repo.ReadGraphEntity(ctx, entityId)
	if err == nil && entityMap != nil {
		// Entity found in Neo4j, extract information
		if kindValue, ok := entityMap["Kind"]; ok {
			kind = &pb.Kind{
				Major: kindValue.(string),
				Minor: "", // Neo4j doesn't store Minor
			}
		}

		if nameValue, ok := entityMap["Name"]; ok {
			// Create a TimeBasedValue with string value
			value, _ := anypb.New(&wrapperspb.StringValue{
				Value: nameValue.(string),
			})

			name = &pb.TimeBasedValue{
				StartTime: entityMap["Created"].(string),
				Value:     value,
			}

			// Add EndTime if available
			if termValue, ok := entityMap["Terminated"]; ok {
				name.EndTime = termValue.(string)
			}
		}

		if createdValue, ok := entityMap["Created"]; ok {
			created = createdValue.(string)
		}

		if termValue, ok := entityMap["Terminated"]; ok {
			terminated = termValue.(string)
		}
	}

	return kind, name, created, terminated, err
}

// GetGraphRelationships retrieves relationships for an entity from Neo4j
func (repo *Neo4jRepository) GetGraphRelationships(ctx context.Context, entityId string) (map[string]*pb.Relationship, error) {
	relationships := make(map[string]*pb.Relationship)
	// Retrieve relationships from Neo4j
	relData, err := repo.ReadRelationships(ctx, entityId)
	if err != nil {
		log.Printf("Error reading relationships for entity %s: %v", entityId, err)
		return relationships, fmt.Errorf("error reading relationships: %v", err)
	}

	// Process each relationship
	for _, rel := range relData {
		relType, ok1 := rel["type"].(string)
		relatedID, ok2 := rel["relatedID"].(string)
		direction, ok3 := rel["direction"].(string)
		created, ok4 := rel["Created"].(string)
		relID, ok5 := rel["relationshipID"].(string)

		if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
			continue // Skip if any required field is missing
		}

		// Create relationship object
		relationship := &pb.Relationship{
			Id:              relID,
			Name:            relType,
			RelatedEntityId: relatedID,
			StartTime:       created,
		}

		// Add termination date if available
		if terminated, ok := rel["Terminated"].(string); ok && terminated != "" {
			relationship.EndTime = terminated
		}

		// Store in map with unique key
		key := fmt.Sprintf("%s-%s-%s", relType, relatedID, direction)
		relationships[key] = relationship
	}

	return relationships, nil
}

// validateGraphEntityCreation checks if an entity has all required fields for Neo4j storage
func validateGraphEntityCreation(entity *pb.Entity) bool {
	// Check if Kind is present and has a Major value
	if entity.Kind == nil || entity.Kind.GetMajor() == "" {
		log.Printf("Skipping Neo4j entity creation for %s: Missing or empty Kind.Major", entity.Id)
		return false
	}

	// Check if Name is present and has a Value
	if entity.Name == nil || entity.Name.GetValue() == nil {
		log.Printf("Skipping Neo4j entity creation for %s: Missing or empty Name.Value", entity.Id)
		return false
	}

	// Check if Created date is present
	if entity.Created == "" {
		log.Printf("Skipping Neo4j entity creation for %s: Missing Created date", entity.Id)
		return false
	}

	return true
}

// Create Graph Entity in Neo4j
func (repo *Neo4jRepository) HandleGraphEntity(ctx context.Context, entity *pb.Entity) error {
	// Validate required fields for Neo4j entity creation
	if validateGraphEntityCreation(entity) {
		// Prepare data for Neo4j with safety checks
		entityMap := map[string]interface{}{
			"Id": entity.Id,
		}

		// Handle Kind field safely
		if entity.Kind != nil {
			entityMap["Kind"] = entity.Kind.GetMajor()
		}

		// Handle Name field safely
		if entity.Name != nil && entity.Name.GetValue() != nil {
			entityMap["Name"] = entity.Name.GetValue().String()
		}

		// Handle other fields
		if entity.Created != "" {
			entityMap["Created"] = entity.Created
		}

		if entity.Terminated != "" {
			entityMap["Terminated"] = entity.Terminated
		}

		// Save entity in Neo4j only if validation passes
		_, err := repo.CreateGraphEntity(ctx, entityMap)
		if err != nil {
			log.Printf("Error saving entity in Neo4j: %v", err)
			// Don't return error here so MongoDB storage is still successful
		} else {
			log.Printf("Successfully saved entity in Neo4j for entity: %s", entity.Id)
		}
	} else {
		log.Printf("Entity %s saved in MongoDB only, skipping Neo4j due to missing required fields", entity.Id)
	}

	return nil
}
