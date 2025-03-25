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
		if majorKindValue, ok := entityMap["MajorKind"]; ok {
			kind = &pb.Kind{
				Major: majorKindValue.(string),
			}
		}

		if minorKindValue, ok := entityMap["MinorKind"]; ok {
			if kind == nil {
				kind = &pb.Kind{}
			}
			kind.Minor = minorKindValue.(string)
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
		log.Printf("[neo4j_handler.GetGraphRelationships] Error reading relationships for entity %s: %v", entityId, err)
		return relationships, fmt.Errorf("[neo4j_handler.GetGraphRelationships] error reading relationships: %v", err)
	}

	// Process each relationship
	// TODO: Holding relationship and defining the content needs to be
	//  revalidated. Discuss and confirm.
	//  Also build a rule based validation for the relationship content.
	for _, rel := range relData {
		relType, ok1 := rel["type"].(string)
		relatedID, ok2 := rel["relatedID"].(string)
		created, ok3 := rel["Created"].(string)
		relID, ok4 := rel["relationshipID"].(string)

		if !ok1 || !ok2 || !ok3 || !ok4 {
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
		relationships[relType] = relationship
	}

	return relationships, nil
}

// validateGraphEntityCreation checks if an entity has all required fields for Neo4j storage
func validateGraphEntityCreation(entity *pb.Entity) bool {
	// Check if Kind is present and has a Major value
	if entity.Kind == nil || entity.Kind.GetMajor() == "" {
		log.Printf("[neo4j_handler.validateGraphEntityCreation] Skipping Neo4j entity creation for %s: Missing or empty Kind.Major", entity.Id)
		return false
	}

	// Check if Name is present and has a Value
	if entity.Name == nil || entity.Name.GetValue() == nil {
		log.Printf("[neo4j_handler.validateGraphEntityCreation] Skipping Neo4j entity creation for %s: Missing or empty Name.Value", entity.Id)
		return false
	}

	// Check if Created date is present
	if entity.Created == "" {
		log.Printf("[neo4j_handler.validateGraphEntityCreation] Skipping Neo4j entity creation for %s: Missing Created date", entity.Id)
		return false
	}

	return true
}

// Create Graph Entity in Neo4j
func (repo *Neo4jRepository) HandleGraphEntity(ctx context.Context, entity *pb.Entity) (bool, error) {
	// Validate required fields for Neo4j entity creation
	if !validateGraphEntityCreation(entity) {
		log.Printf("[neo4j_handler.HandleGraphEntity] Entity %s saved in MongoDB only, skipping Neo4j due to missing required fields", entity.Id)
		return false, fmt.Errorf("[neo4j_handler.HandleGraphEntity] missing required fields for Neo4j entity creation")
	}

	// Prepare data for Neo4j with safety checks
	entityMap := map[string]interface{}{
		"Id": entity.Id,
	}

	// Validate and extract the Kind field
	if entity.Kind == nil || entity.Kind.GetMajor() == "" || entity.Kind.GetMinor() == "" {
		return false, fmt.Errorf("[neo4j_handler.HandleGraphEntity] missing or invalid Kind.Major or Kind.Minor for entity %s", entity.Id)
	}

	kind := &pb.Kind{
		Major: entity.Kind.GetMajor(),
		Minor: entity.Kind.GetMinor(),
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

	// First check if entity exists
	existingEntity, err := repo.ReadGraphEntity(ctx, entity.Id)
	if err == nil && existingEntity != nil {
		// Entity exists, update it
		log.Printf("[neo4j_handler.HandleGraphEntity] Updating existing entity in Neo4j: %s", entity.Id)
		result, err := repo.UpdateGraphEntity(ctx, entity.Id, entityMap)
		if err != nil {
			log.Printf("[neo4j_handler.HandleGraphEntity] Error updating entity in Neo4j: %v", err)
			return false, err
		} else {
			log.Printf("[neo4j_handler.HandleGraphEntity] Successfully updated entity in Neo4j: %s", entity.Id)
			return result != nil, nil // Success if we got a non-nil result
		}
	} else {
		// Entity doesn't exist, create it
		log.Printf("[neo4j_handler.HandleGraphEntity] Creating new entity in Neo4j: %s", entity.Id)
		result, err := repo.CreateGraphEntity(ctx, kind, entityMap)
		if err != nil {
			log.Printf("[neo4j_handler.HandleGraphEntity] Error creating entity in Neo4j: %v", err)
			return false, err
		} else {
			log.Printf("[neo4j_handler.HandleGraphEntity] Successfully created entity in Neo4j: %s", entity.Id)
			return result != nil, nil // Success if we got a non-nil result
		}
	}
}

// Handle Graph Relationships in Neo4j
func (repo *Neo4jRepository) HandleGraphRelationships(ctx context.Context, entity *pb.Entity) error {
	// Check if entity has relationships
	if len(entity.Relationships) == 0 {
		log.Printf("[neo4j_handler.HandleGraphRelationships] No relationships to process for entity: %s", entity.Id)
		return nil
	}

	log.Printf("[neo4j_handler.HandleGraphRelationships] Processing %d relationships for entity: %s", len(entity.Relationships), entity.Id)

	// First, process all child entities
	for _, relationship := range entity.Relationships {
		if relationship == nil || relationship.RelatedEntityId == "" {
			continue
		}

		// Check if the child entity exists
		childEntityMap, err := repo.ReadGraphEntity(ctx, relationship.RelatedEntityId)
		if err != nil || childEntityMap == nil {
			// Child entity doesn't exist, log this information
			log.Printf("[neo4j_handler.HandleGraphRelationships] Child entity %s does not exist in Neo4j. Make sure to create it first.",
				relationship.RelatedEntityId)
			return fmt.Errorf("[neo4j_handler.HandleGraphRelationships] child entity %s does not exist", relationship.RelatedEntityId)
		}
		log.Printf("[neo4j_handler.HandleGraphRelationships] Child entity %s exists in Neo4j", relationship.RelatedEntityId)
	}

	// Now that child entities are handled, create or update relationships from parent to children
	for relationshipType, relationship := range entity.Relationships {
		if relationship == nil {
			continue
		}

		// Skip if related entity ID is empty
		if relationship.RelatedEntityId == "" {
			log.Printf("[neo4j_handler.HandleGraphRelationships] Skipping relationship with empty related entity ID for entity: %s", entity.Id)
			continue
		}

		// Check if this relationship already exists (if it has an ID)
		if relationship.Id != "" {
			// Prepare relationship data for update
			relationshipData := map[string]interface{}{
				"relatedEntityId": relationship.RelatedEntityId,
				"name":            relationship.Name,
				"startTime":       relationship.StartTime,
			}

			if relationship.EndTime != "" {
				relationshipData["endTime"] = relationship.EndTime
			}

			// Try to update the relationship
			_, err := repo.UpdateRelationship(ctx, relationship.Id, relationshipData)
			if err != nil {
				log.Printf("Error updating relationship %s from %s to %s: %v",
					relationship.Id, entity.Id, relationship.RelatedEntityId, err)

				// If update fails (perhaps relationship doesn't exist despite having ID),
				// try to create it instead
				_, createErr := repo.CreateRelationship(ctx, entity.Id, relationship)
				if createErr != nil {
					log.Printf("[neo4j_handler.HandleGraphRelationships] Also failed to create relationship: %v", createErr)
				} else {
					log.Printf("[neo4j_handler.HandleGraphRelationships] Successfully created relationship %s after update failure", relationship.Id)
				}
			} else {
				log.Printf("[neo4j_handler.HandleGraphRelationships] Successfully updated relationship %s from %s to %s of type %s",
					relationship.Id, entity.Id, relationship.RelatedEntityId, relationshipType)
			}
		} else {
			// No ID, so create a new relationship
			_, err := repo.CreateRelationship(ctx, entity.Id, relationship)
			if err != nil {
				log.Printf("[neo4j_handler.HandleGraphRelationships] Error creating relationship from %s to %s of type %s: %v",
					entity.Id, relationship.RelatedEntityId, relationshipType, err)
			} else {
				log.Printf("[neo4j_handler.HandleGraphRelationships] Successfully created relationship from %s to %s of type %s",
					entity.Id, relationship.RelatedEntityId, relationshipType)
			}
		}
	}

	return nil
}
