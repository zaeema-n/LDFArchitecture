package neo4jrepository

import (
	"context"
	"fmt"
	"lk/datafoundation/crud-api/db/config"
	pb "lk/datafoundation/crud-api/lk/datafoundation/crud-api"
	"log"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jRepository struct {
	client neo4j.DriverWithContext
	config *config.Neo4jConfig
}

// NewNeo4jRepository initializes a Neo4j driver
func NewNeo4jRepository(ctx context.Context, config *config.Neo4jConfig) (*Neo4jRepository, error) {
	client, err := neo4j.NewDriverWithContext(config.URI, neo4j.BasicAuth(config.Username, config.Password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j driver: %w", err)
	}

	// Verify connectivity
	if err := client.VerifyConnectivity(ctx); err != nil {
		client.Close(ctx) // Close if connectivity check fails
		return nil, fmt.Errorf("failed to connect to Neo4j: %w", err)
	}

	log.Println("Connected to Neo4j successfully!")

	return &Neo4jRepository{
		client: client,
		config: config,
	}, nil
}

// Close properly closes the Neo4j driver
func (r *Neo4jRepository) Close(ctx context.Context) {
	if r.client != nil {
		r.client.Close(ctx)
		log.Println("Neo4j connection closed")
	}
}

// getSession creates a new session
func (r *Neo4jRepository) getSession(ctx context.Context) neo4j.SessionWithContext {
	return r.client.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
}

// createEntity checks if an entity exists and creates it if it doesn’t
func (r *Neo4jRepository) createEntity(ctx context.Context, entityMap map[string]interface{}) (map[string]interface{}, error) {
	// Ensure the map has the necessary fields
	id, ok := entityMap["id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'id' field")
	}

	kind, ok := entityMap["kind"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'kind' field")
	}

	name, ok := entityMap["name"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'name' field")
	}

	dateCreated, ok := entityMap["dateCreated"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'dateCreated' field")
	}

	// Optional field
	var dateEnded *string
	if ended, ok := entityMap["dateEnded"].(string); ok {
		dateEnded = &ended
	}

	// Convert ID to integer (or string if needed)
	entityID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Check if the node already exists
	existsQuery := `MATCH (e:` + kind + ` {id: $id}) RETURN e`
	result, err := session.Run(ctx, existsQuery, map[string]interface{}{"id": entityID})
	if err != nil {
		return nil, fmt.Errorf("error checking if entity exists: %v", err)
	}

	// If entity exists, return an error
	if result.Next(ctx) {
		return nil, fmt.Errorf("entity with ID %s already exists", id)
	}

	// Create the node
	createQuery := `CREATE (e:` + kind + ` {id: $id, name: $name, dateCreated: date($dateCreated)`
	if dateEnded != nil {
		createQuery += `, dateEnded: date($dateEnded)`
	}
	createQuery += `}) RETURN e`

	// Set parameters for the query
	params := map[string]interface{}{
		"id":          entityID,
		"name":        name,
		"dateCreated": dateCreated,
	}
	if dateEnded != nil {
		params["dateEnded"] = *dateEnded
	}

	// Run the query to create the entity and return it
	result, err = session.Run(ctx, createQuery, params)
	if err != nil {
		return nil, fmt.Errorf("error creating entity: %v", err)
	}

	// Retrieve the created entity (assuming it returns a node)
	if result.Next(ctx) {
		createdEntity, _ := result.Record().Get("e")

		// Convert the node to a map
		node, ok := createdEntity.(neo4j.Node)
		if !ok {
			return nil, fmt.Errorf("failed to cast created entity to neo4j.Node")
		}

		// Convert all properties to strings
		createdEntityMap := map[string]interface{}{
			"id":          fmt.Sprintf("%v", node.Props["id"]),          // Convert to string
			"name":        fmt.Sprintf("%v", node.Props["name"]),        // Convert to string
			"dateCreated": fmt.Sprintf("%v", node.Props["dateCreated"]), // Convert to string
		}
		if dateEnded != nil {
			createdEntityMap["dateEnded"] = fmt.Sprintf("%v", *dateEnded) // Convert to string
		}

		// Return the created entity
		return createdEntityMap, nil
	}

	return nil, fmt.Errorf("failed to retrieve created entity")
}

func (r *Neo4jRepository) createRelationship(ctx context.Context, entityID string, rel *pb.Relationship) (map[string]interface{}, error) {
	// Convert entityID (parent ID) to integer
	parentID, err := strconv.Atoi(entityID)
	if err != nil {
		return nil, fmt.Errorf("invalid parent ID format: %v", err)
	}

	// Convert relatedEntityId (child ID) to integer
	childID, err := strconv.Atoi(rel.RelatedEntityId)
	if err != nil {
		return nil, fmt.Errorf("invalid child ID format: %v", err)
	}

	// Convert relationship ID to integer
	relationshipID, err := strconv.Atoi(rel.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid relationship ID format: %v", err)
	}

	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Check if both entities exist
	existsQuery := `MATCH (p {id: $parentID}), (c {id: $childID}) RETURN p, c`
	result, err := session.Run(ctx, existsQuery, map[string]interface{}{
		"parentID": parentID,
		"childID":  childID,
	})
	if err != nil {
		return nil, fmt.Errorf("error checking entities: %v", err)
	}
	if !result.Next(ctx) {
		return nil, fmt.Errorf("either parent or child entity does not exist")
	}

	// Build the Cypher query to create the relationship
	createQuery := `MATCH (p {id: $parentID}), (c {id: $childID})
                    MERGE (p)-[r:` + rel.Name + ` {id: $relationshipID}]->(c)
                    SET r.startDate = date($startDate)
                    RETURN r`

	// Parameters for the query
	params := map[string]interface{}{
		"parentID":       parentID,
		"childID":        childID,
		"relationshipID": relationshipID,
		"startDate":      rel.StartTime,
	}

	// If EndTime exists, add it to the query and parameters
	if rel.EndTime != "" {
		createQuery += `, r.endDate = date($endDate)`
		params["endDate"] = rel.EndTime
	}

	// Run the query to create the relationship and return it
	result, err = session.Run(ctx, createQuery, params)
	if err != nil {
		return nil, fmt.Errorf("error creating relationship: %v", err)
	}

	// Retrieve the created relationship
	if result.Next(ctx) {
		createdRel, _ := result.Record().Get("r")

		// Convert the relationship to a map
		relationship, ok := createdRel.(neo4j.Relationship)
		if !ok {
			return nil, fmt.Errorf("failed to cast created relationship to neo4j.Relationship")
		}

		// Convert all properties to strings
		relationshipMap := map[string]interface{}{
			"id":               fmt.Sprintf("%v", relationship.Props["id"]),
			"startDate":        fmt.Sprintf("%v", relationship.Props["startDate"]),
			"relationshipType": rel.Name,
		}
		if rel.EndTime != "" {
			relationshipMap["endDate"] = fmt.Sprintf("%v", relationship.Props["endDate"])
		}

		// Return the created relationship
		return relationshipMap, nil
	}

	return nil, fmt.Errorf("failed to retrieve created relationship")
}

// ReadEntity retrieves an entity by its ID from the Neo4j database and returns it as a map.
func (r *Neo4jRepository) ReadEntity(ctx context.Context, entityID string) (map[string]interface{}, error) {
	// Ensure the entity ID is valid
	if entityID == "" {
		return nil, fmt.Errorf("entity ID cannot be empty")
	}

	// Convert entity ID to an integer
	entityInt, err := strconv.Atoi(entityID)
	if err != nil {
		return nil, fmt.Errorf("invalid entity ID format: %v", err)
	}

	// Get a session
	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Query to retrieve entity details
	query := `MATCH (e {id: $id})
              RETURN labels(e)[0] AS kind, e.id AS id, e.name AS name, 
                     toString(e.dateCreated) AS dateCreated, 
                     CASE WHEN e.dateEnded IS NOT NULL THEN toString(e.dateEnded) ELSE NULL END AS dateEnded`

	// Execute the query
	result, err := session.Run(ctx, query, map[string]interface{}{"id": entityInt})
	if err != nil {
		return nil, fmt.Errorf("error querying entity: %v", err)
	}

	// Extract and process entity data
	if result.Next(ctx) {
		record := result.Record()

		entity := map[string]interface{}{
			"id":          fmt.Sprintf("%v", record.Values[1]), // Convert ID back to string
			"kind":        fmt.Sprintf("%v", record.Values[0]),
			"name":        fmt.Sprintf("%v", record.Values[2]),
			"dateCreated": fmt.Sprintf("%v", record.Values[3]),
		}

		// Handle optional dateEnded field
		if dateEndedVal, exists := record.Get("dateEnded"); exists && dateEndedVal != nil {
			entity["dateEnded"] = fmt.Sprintf("%v", dateEndedVal)
		}

		return entity, nil
	}

	return nil, fmt.Errorf("entity with ID %s not found", entityID)
}

// ReadRelatedEntityIds retrieves related entity IDs based on a given relationship
func (r *Neo4jRepository) ReadRelatedEntityIds(ctx context.Context, entityID string, relationship string, ts string) ([]string, error) {
	// Ensure the entity ID is valid
	if entityID == "" {
		return nil, fmt.Errorf("entity ID cannot be empty")
	}

	entityInt, err := strconv.Atoi(entityID)
	if err != nil {
		return nil, fmt.Errorf("invalid entity ID format: %v", err)
	}

	// Open session
	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Query to find related entities within the given timestamp
	query := fmt.Sprintf(`
		MATCH (e {id: $entityID})-[r:%s]->(related)
		WHERE r.startDate <= date($ts) AND (r.endDate IS NULL OR r.endDate > date($ts))
		RETURN related.id AS relatedID
	`, relationship)

	// Run the query
	result, err := session.Run(ctx, query, map[string]interface{}{
		"entityID": entityInt,
		"ts":       ts, // Date format should be "YYYY-MM-DD"
	})
	if err != nil {
		return nil, fmt.Errorf("error querying related entities: %v", err)
	}

	// Extract related entity IDs
	var relatedIDs []string
	for result.Next(ctx) {
		record := result.Record()
		if relatedID, exists := record.Get("relatedID"); exists && relatedID != nil {
			relatedIDs = append(relatedIDs, fmt.Sprintf("%v", relatedID))
		}
	}

	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over query result: %v", err)
	}

	return relatedIDs, nil
}

func (r *Neo4jRepository) readRelationships(ctx context.Context, entityID string) ([]map[string]interface{}, error) {
	// Convert entity ID to integer
	entityInt, err := strconv.Atoi(entityID)
	if err != nil {
		return nil, fmt.Errorf("invalid entity ID format: %v", err)
	}

	// Open session
	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Cypher query to get all relationships (incoming and outgoing)
	query := `
		MATCH (e {id: $entityID})-[r]->(related)
		RETURN type(r) AS type, related.id AS relatedID, "OUTGOING" AS direction, 
		       toString(r.startDate) AS startDate, 
		       CASE WHEN r.endDate IS NOT NULL THEN toString(r.endDate) ELSE NULL END AS endDate,
		       r.id AS relationshipID
		UNION
		MATCH (e {id: $entityID})<-[r]-(related)
		RETURN type(r) AS type, related.id AS relatedID, "INCOMING" AS direction, 
		       toString(r.startDate) AS startDate, 
		       CASE WHEN r.endDate IS NOT NULL THEN toString(r.endDate) ELSE NULL END AS endDate,
		       r.id AS relationshipID
	`

	// Run the query
	result, err := session.Run(ctx, query, map[string]interface{}{
		"entityID": entityInt,
	})
	if err != nil {
		return nil, fmt.Errorf("error querying relationships: %v", err)
	}

	// Process results
	var relationships []map[string]interface{}
	for result.Next(ctx) {
		record := result.Record()
		values := record.Values

		// Ensure expected values exist
		if len(values) < 6 {
			continue
		}

		// Relationship structure
		rel := map[string]interface{}{
			"type":           fmt.Sprintf("%v", values[0]), // Relationship type
			"relatedID":      fmt.Sprintf("%v", values[1]),
			"direction":      fmt.Sprintf("%v", values[2]), // "INCOMING" or "OUTGOING"
			"startDate":      fmt.Sprintf("%v", values[3]),
			"relationshipID": fmt.Sprintf("%v", values[5]), // Relationship ID
		}

		// Optional endDate
		if values[4] != nil {
			rel["endDate"] = fmt.Sprintf("%v", values[4])
		}

		relationships = append(relationships, rel)
	}

	// Return relationships as a map
	return relationships, nil
}

func (r *Neo4jRepository) getRelationship(ctx context.Context, relationshipID string) (map[string]interface{}, error) {
	// Convert relationship ID to integer
	relID, err := strconv.Atoi(relationshipID)
	if err != nil {
		return nil, fmt.Errorf("invalid relationship ID format: %v", err)
	}

	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Cypher query to find the relationship by its ID
	query := `
		MATCH ()-[r]->()
		WHERE r.id = $relationshipID
		RETURN type(r) AS type, startNode(r).id AS startEntityID, endNode(r).id AS endEntityID, 
		       toString(r.startDate) AS startDate, 
		       CASE WHEN r.endDate IS NOT NULL THEN toString(r.endDate) ELSE NULL END AS endDate, 
		       r.id AS relationshipID
	`

	// Run the query to fetch the relationship
	result, err := session.Run(ctx, query, map[string]interface{}{
		"relationshipID": relID,
	})
	if err != nil {
		return nil, fmt.Errorf("error querying relationship: %v", err)
	}

	// Process results
	if result.Next(ctx) {
		record := result.Record()
		values := record.Values

		// Ensure expected values exist
		if len(values) < 6 {
			return nil, fmt.Errorf("unexpected data format for relationship")
		}

		// Map to hold the relationship data
		relationship := map[string]interface{}{
			"type":           fmt.Sprintf("%v", values[0]), // Relationship type
			"startEntityID":  fmt.Sprintf("%v", values[1]), // ID of the start entity
			"endEntityID":    fmt.Sprintf("%v", values[2]), // ID of the end entity
			"startDate":      fmt.Sprintf("%v", values[3]),
			"relationshipID": fmt.Sprintf("%v", values[5]), // Relationship ID
		}

		// Optional endDate
		if values[4] != nil {
			relationship["endDate"] = fmt.Sprintf("%v", values[4])
		}

		// Return the relationship data as a map
		return relationship, nil
	}

	// If no relationship was found
	return nil, fmt.Errorf("relationship with ID %s not found", relationshipID)
}

// func (r *Neo4jRepository) updateEntity(ctx context.Context, entityID string, updateData map[string]interface{}) (map[string]interface{}, error) {
// 	// Convert entity ID to integer
// 	entityInt, err := strconv.Atoi(entityID)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid entity ID format: %v", err)
// 	}

// 	// Open session
// 	session := r.getSession(ctx)
// 	defer session.Close(ctx)

// 	// Check if the entity exists
// 	existsQuery := `MATCH (e) WHERE e.id = $id RETURN e`
// 	result, err := session.Run(ctx, existsQuery, map[string]interface{}{"id": entityInt})
// 	if err != nil {
// 		return nil, fmt.Errorf("error checking if entity exists: %v", err)
// 	}

// 	// If entity doesn't exist, return error
// 	if !result.Next(ctx) {
// 		return nil, fmt.Errorf("entity with ID %d does not exist", entityInt)
// 	}

// 	// Prepare the query and parameters dynamically
// 	updateQuery := `MATCH (e) WHERE e.id = $id `
// 	params := map[string]interface{}{"id": entityInt}

// 	// Dynamically build SET clause for all fields in updateData
// 	setClauses := []string{}
// 	for key, value := range updateData {
// 		params[key] = value
// 		setClauses = append(setClauses, fmt.Sprintf("e.%s = $%s", key, key))
// 	}

// 	// Join all SET clauses with commas
// 	if len(setClauses) > 0 {
// 		updateQuery += "SET " + strings.Join(setClauses, ", ")
// 	}

// 	// Execute the update query
// 	_, err = session.Run(ctx, updateQuery, params)
// 	if err != nil {
// 		return nil, fmt.Errorf("error updating entity: %v", err)
// 	}

// 	// Fetch the updated entity
// 	updatedEntityQuery := `MATCH (e) WHERE e.id = $id RETURN e`
// 	result, err = session.Run(ctx, updatedEntityQuery, map[string]interface{}{"id": entityInt})
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching updated entity: %v", err)
// 	}

// 	// Retrieve the updated entity
// 	if result.Next(ctx) {
// 		updatedEntity, _ := result.Record().Get("e")

// 		// Convert the node to a map
// 		node, ok := updatedEntity.(neo4j.Node)
// 		if !ok {
// 			return nil, fmt.Errorf("failed to cast updated entity to neo4j.Node")
// 		}

// 		// Dynamically convert all properties to strings
// 		updatedEntityMap := make(map[string]interface{})
// 		for key, value := range node.Props {
// 			updatedEntityMap[key] = fmt.Sprintf("%v", value) // Convert each property to a string
// 		}

// 		// Return the updated entity
// 		return updatedEntityMap, nil
// 	}

// 	return nil, fmt.Errorf("failed to retrieve updated entity")
// }

// this function can only be used to update the name and dateEnded fields of an entity
func (r *Neo4jRepository) updateEntity(ctx context.Context, id string, updateData map[string]interface{}) (map[string]interface{}, error) {
	// Convert entity ID to integer
	entityInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid entity ID format: %v", err)
	}

	// Prepare update parameters
	params := map[string]interface{}{
		"id": entityInt,
	}

	// Open session
	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Check if the entity exists
	existsQuery := `MATCH (e {id: $id}) RETURN e`
	result, err := session.Run(ctx, existsQuery, params)
	if err != nil {
		return nil, fmt.Errorf("error checking if entity exists: %v", err)
	}

	if !result.Next(ctx) {
		return nil, fmt.Errorf("entity with ID %d does not exist", entityInt)
	}

	// Build Cypher query for updating entity
	query := `
        MATCH (e {id: $id})
    `

	// Add `name` if provided
	if name, exists := updateData["name"]; exists {
		params["name"] = name
		query += `SET e.name = $name `
	}

	// Add `dateEnded` if provided
	if dateEnded, exists := updateData["dateEnded"]; exists {
		params["dateEnded"] = dateEnded
		query += `SET e.dateEnded = date($dateEnded) `
	}

	// Execute update query and return updated entity
	query += ` RETURN e`

	result, err = session.Run(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("error updating entity: %v", err)
	}

	// Retrieve updated entity
	if result.Next(ctx) {
		node, ok := result.Record().Get("e")
		if !ok {
			return nil, fmt.Errorf("unexpected error retrieving entity")
		}

		// Convert node properties to map
		entityNode := node.(neo4j.Node)
		updatedEntity := make(map[string]interface{})
		for key, value := range entityNode.Props {
			updatedEntity[key] = fmt.Sprintf("%v", value)
		}

		return updatedEntity, nil
	}

	return nil, fmt.Errorf("failed to retrieve updated entity")
}

func (r *Neo4jRepository) updateRelationship(ctx context.Context, relationshipID string, updateData map[string]interface{}) (map[string]interface{}, error) {
	// Convert relationship ID to integer
	relationshipInt, err := strconv.Atoi(relationshipID)
	if err != nil {
		return nil, fmt.Errorf("invalid relationship ID format: %v", err)
	}

	// Prepare update parameters
	params := map[string]interface{}{
		"relationshipID": relationshipInt,
	}

	// Open session
	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Check if the relationship exists
	existsQuery := `MATCH ()-[r {id: $relationshipID}]->() RETURN r`
	result, err := session.Run(ctx, existsQuery, params)
	if err != nil {
		return nil, fmt.Errorf("error checking if relationship exists: %v", err)
	}

	if !result.Next(ctx) {
		return nil, fmt.Errorf("relationship with ID %d does not exist", relationshipInt)
	}

	// Build Cypher query for updating relationship
	query := `
        MATCH ()-[r {id: $relationshipID}]->()
    `

	// Add `endDate` if provided (required)
	endDate, exists := updateData["endDate"]
	if !exists {
		return nil, fmt.Errorf("endDate is required")
	}
	params["endDate"] = endDate
	query += `SET r.endDate = date($endDate) RETURN r`

	// Execute update query and return updated relationship
	result, err = session.Run(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("error updating relationship: %v", err)
	}

	// Retrieve updated relationship
	if result.Next(ctx) {
		rel, ok := result.Record().Get("r")
		if !ok {
			return nil, fmt.Errorf("unexpected error retrieving relationship")
		}

		// Convert relationship properties to map with string values
		relationship := rel.(neo4j.Relationship)
		updatedRelationship := make(map[string]interface{})
		for key, value := range relationship.Props {
			updatedRelationship[key] = fmt.Sprintf("%v", value) // Convert everything to string
		}

		return updatedRelationship, nil
	}

	return nil, fmt.Errorf("failed to retrieve updated relationship")
}

func (r *Neo4jRepository) deleteRelationship(ctx context.Context, relationshipID string) error {
	// Convert relationship ID to integer
	relationshipInt, err := strconv.Atoi(relationshipID)
	if err != nil {
		return fmt.Errorf("invalid relationship ID format: %v", err)
	}

	// Prepare query parameters
	params := map[string]interface{}{
		"relationshipID": relationshipInt,
	}

	// Open session
	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Check if the relationship exists
	query := `MATCH ()-[r {id: $relationshipID}]->() RETURN r`
	result, err := session.Run(ctx, query, params)
	if err != nil {
		return fmt.Errorf("error checking if relationship exists: %v", err)
	}

	// If no relationship is found, return an error
	if !result.Next(ctx) {
		return fmt.Errorf("relationship with ID %d does not exist", relationshipInt)
	}

	// Delete the relationship
	deleteQuery := `MATCH ()-[r {id: $relationshipID}]->() DELETE r`
	_, err = session.Run(ctx, deleteQuery, params)
	if err != nil {
		return fmt.Errorf("error deleting relationship: %v", err)
	}

	return nil
}

func (r *Neo4jRepository) deleteEntity(ctx context.Context, entityID string) error {
	// Convert entity ID to integer
	entityInt, err := strconv.Atoi(entityID)
	if err != nil {
		return fmt.Errorf("invalid entity ID format: %v", err)
	}

	// Open a session
	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Check if the entity exists before attempting to delete
	query := `MATCH (e {id: $entityID}) RETURN e`
	params := map[string]interface{}{
		"entityID": entityInt,
	}

	// Run query to find the entity
	result, err := session.Run(ctx, query, params)
	if err != nil {
		return fmt.Errorf("error checking if entity exists: %v", err)
	}

	// If entity doesn't exist, return an error
	if !result.Next(ctx) {
		return fmt.Errorf("entity with ID %d does not exist", entityInt)
	}

	// Get the relationships of the entity
	relationships, err := r.readRelationships(ctx, entityID)
	if err != nil {
		return fmt.Errorf("error getting relationships: %v", err)
	}

	// If there are relationships, return an error with relationship details
	if len(relationships) > 0 {
		return fmt.Errorf("entity has relationships and cannot be deleted. Relationships: %v", relationships)
	}

	// Delete the entity (node) with the given ID
	deleteQuery := `MATCH (e {id: $entityID}) DELETE e`
	_, err = session.Run(ctx, deleteQuery, params)
	if err != nil {
		return fmt.Errorf("error deleting entity: %v", err)
	}

	return nil
}




