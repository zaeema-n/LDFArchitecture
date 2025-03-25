package neo4jrepository

import (
	"context"
	"fmt"
	"lk/datafoundation/crud-api/db/config"
	pb "lk/datafoundation/crud-api/lk/datafoundation/crud-api"
	"log"

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
		log.Printf("[neo4j_client.NewNeo4jRepository] failed to create Neo4j driver: %v", err)
		return nil, fmt.Errorf("failed to create Neo4j driver: %w", err)
	}

	// Verify connectivity
	if err := client.VerifyConnectivity(ctx); err != nil {
		client.Close(ctx) // Close if connectivity check fails
		log.Printf("[neo4j_client.NewNeo4jRepository] failed to connect to Neo4j: %v", err)
		return nil, fmt.Errorf("failed to connect to Neo4j: %w", err)
	}

	log.Println("[neo4j_client.NewNeo4jRepository] Connected to Neo4j successfully!")

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

// CreateGraphEntity checks if an entity exists and creates it if it doesn't
func (r *Neo4jRepository) CreateGraphEntity(ctx context.Context, kind *pb.Kind, entityMap map[string]interface{}) (map[string]interface{}, error) {
	// Validate the kind parameter
	if kind == nil || kind.Major == "" {
		log.Printf("[neo4j_client.CreateGraphEntity] missing or invalid 'Kind.Major' field")
		return nil, fmt.Errorf("[neo4j_client.CreateGraphEntity] missing or invalid 'Kind.Major' field")
	} else {
		log.Printf("[neo4j_client.CreateGraphEntity] Kind.Major: %v", kind.Major)
	}

	// Extract the required fields from the entityMap
	id, ok := entityMap["Id"].(string)
	if !ok {
		log.Printf("[neo4j_client.CreateGraphEntity] missing or invalid 'Id' field")
		return nil, fmt.Errorf("[neo4j_client.CreateGraphEntity] missing or invalid 'Id' field")
	} else {
		log.Printf("[neo4j_client.CreateGraphEntity] Id: %v", id)
	}

	name, ok := entityMap["Name"].(string)
	if !ok {
		log.Printf("[neo4j_client.CreateGraphEntity] missing or invalid 'Name' field")
		return nil, fmt.Errorf("[neo4j_client.CreateGraphEntity] missing or invalid 'Name' field")
	} else {
		log.Printf("[neo4j_client.CreateGraphEntity] Name: %v", name)
	}

	created, ok := entityMap["Created"].(string)
	if !ok {
		log.Printf("[neo4j_client.CreateGraphEntity] missing or invalid 'Created' field")
		return nil, fmt.Errorf("[neo4j_client.CreateGraphEntity] missing or invalid 'Created' field")
	} else {
		log.Printf("[neo4j_client.CreateGraphEntity] Created: %v", created)
	}

	// Optional field
	var terminated *string
	if term, ok := entityMap["Terminated"].(string); ok {
		terminated = &term
	} else {
		log.Printf("[neo4j_client.CreateGraphEntity] Terminated: %v", terminated)
	}

	// Open a session
	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Check if the node already exists
	existsQuery := `MATCH (e:` + kind.Major + ` {Id: $Id}) RETURN e`
	result, err := session.Run(ctx, existsQuery, map[string]interface{}{"Id": id})
	if err != nil {
		log.Printf("[neo4j_client.CreateGraphEntity] error checking if entity exists: %v", err)
		return nil, fmt.Errorf("[neo4j_client.CreateGraphEntity] error checking if entity exists: %v", err)
	} else {
		log.Printf("[neo4j_client.CreateGraphEntity] existsQuery: %v", existsQuery)
	}

	// If entity exists, return an error
	if result.Next(ctx) {
		log.Printf("[neo4j_client.CreateGraphEntity] entity with Id %s already exists", id)
		return nil, fmt.Errorf("[neo4j_client.CreateGraphEntity] entity with Id %s already exists", id)
	} else {
		log.Printf("[neo4j_client.CreateGraphEntity] entity with Id %s does not exist", id)
	}

	// Create the node
	createQuery := `CREATE (e:` + kind.Major + ` {Id: $Id, Name: $Name, Created: date($Created), MinorKind: $MinorKind`
	if terminated != nil {
		createQuery += `, Terminated: date($Terminated)`
	}
	createQuery += `}) RETURN e`

	// Set parameters for the query
	params := map[string]interface{}{
		"Id":        id,
		"Name":      name,
		"Created":   created,
		"MinorKind": kind.Minor,
	}
	if terminated != nil {
		params["Terminated"] = *terminated
	}

	// Run the query to create the entity and return it
	result, err = session.Run(ctx, createQuery, params)
	if err != nil {
		log.Printf("[neo4j_client.CreateGraphEntity] error creating entity: %v", err)
		return nil, fmt.Errorf("[neo4j_client.CreateGraphEntity] error creating entity: %v", err)
	} else {
		log.Printf("[neo4j_client.CreateGraphEntity] created entity: %v", params)
	}

	// Retrieve the created entity
	if result.Next(ctx) {
		createdEntity, _ := result.Record().Get("e")
		node, ok := createdEntity.(neo4j.Node)
		if !ok {
			log.Printf("[neo4j_client.CreateGraphEntity] failed to cast created entity to neo4j.Node")
			return nil, fmt.Errorf("[neo4j_client.CreateGraphEntity] failed to cast created entity to neo4j.Node")
		} else {
			log.Printf("[neo4j_client.CreateGraphEntity] created entity: %v", createdEntity)
		}

		// Convert the node properties to a map
		createdEntityMap := map[string]interface{}{
			"Id":        fmt.Sprintf("%v", node.Props["Id"]),
			"Name":      fmt.Sprintf("%v", node.Props["Name"]),
			"Created":   fmt.Sprintf("%v", node.Props["Created"]),
			"MinorKind": fmt.Sprintf("%v", node.Props["MinorKind"]),
		}
		if terminated != nil {
			createdEntityMap["Terminated"] = fmt.Sprintf("%v", *terminated)
		} else {
			log.Printf("[neo4j_client.CreateGraphEntity] Terminated: %v", terminated)
		}
		log.Printf("[neo4j_client.CreateGraphEntity] created entity: %v", createdEntityMap)
		return createdEntityMap, nil
	}

	log.Printf("[neo4j_client.CreateGraphEntity] failed to create entity")
	return nil, fmt.Errorf("[neo4j_client.CreateGraphEntity] failed to create entity")
}

// CreateRelationship creates a relationship between two entities
func (r *Neo4jRepository) CreateRelationship(ctx context.Context, entityID string, rel *pb.Relationship) (map[string]interface{}, error) {
	session := r.getSession(ctx)
	defer session.Close(ctx)

	existsQuery := `MATCH (p {Id: $parentID}), (c {Id: $childID}) RETURN p, c`
	result, err := session.Run(ctx, existsQuery, map[string]interface{}{
		"parentID": entityID,
		"childID":  rel.RelatedEntityId,
	})
	if err != nil {
		log.Printf("[neo4j_client.CreateRelationship] error checking entities: %v", err)
		return nil, fmt.Errorf("error checking entities: %v", err)
	} else {
		log.Printf("[neo4j_client.CreateRelationship] existsQuery: %v", existsQuery)
	}
	if !result.Next(ctx) {
		log.Printf("[neo4j_client.CreateRelationship] either parent or child entity does not exist")
		return nil, fmt.Errorf("either parent or child entity does not exist")
	} else {
		log.Printf("[neo4j_client.CreateRelationship] either parent or child entity does not exist")
	}

	createQuery := `MATCH (p {Id: $parentID}), (c {Id: $childID})
                    MERGE (p)-[r:` + rel.Name + ` {Id: $relationshipID}]->(c)
                    SET r.Created = date($startDate)
                    RETURN r`

	params := map[string]interface{}{
		"parentID":       entityID,
		"childID":        rel.RelatedEntityId,
		"relationshipID": rel.Id,
		"startDate":      rel.StartTime,
	}

	if rel.EndTime != "" {
		createQuery += `, r.Terminated = date($endDate)`
		params["endDate"] = rel.EndTime
	}

	result, err = session.Run(ctx, createQuery, params)
	if err != nil {
		log.Printf("[neo4j_client.CreateRelationship] error creating relationship: %v", err)
		return nil, fmt.Errorf("error creating relationship: %v", err)
	} else {
		log.Printf("[neo4j_client.CreateRelationship] createQuery: %v", createQuery)
	}

	if result.Next(ctx) {
		createdRel, _ := result.Record().Get("r")
		relationship, ok := createdRel.(neo4j.Relationship)
		if !ok {
			log.Printf("[neo4j_client.CreateRelationship] failed to cast created relationship to neo4j.Relationship")
			return nil, fmt.Errorf("failed to cast created relationship to neo4j.Relationship")
		} else {
			log.Printf("[neo4j_client.CreateRelationship] created relationship: %v", createdRel)
		}

		relationshipMap := map[string]interface{}{
			"Id":               fmt.Sprintf("%v", relationship.Props["Id"]),
			"Created":          fmt.Sprintf("%v", relationship.Props["Created"]),
			"relationshipType": rel.Name,
		}
		if rel.EndTime != "" {
			relationshipMap["Terminated"] = fmt.Sprintf("%v", relationship.Props["Terminated"])
		}
		log.Printf("[neo4j_client.CreateRelationship] created relationship: %v", relationshipMap)
		return relationshipMap, nil
	}

	return nil, fmt.Errorf("failed to retrieve created relationship")
}

// ReadGraphEntity retrieves an entity by its ID from the Neo4j database and returns it as a map.
func (r *Neo4jRepository) ReadGraphEntity(ctx context.Context, entityID string) (map[string]interface{}, error) {
	if entityID == "" {
		return nil, fmt.Errorf("entity Id cannot be empty")
	}

	// Open a session
	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Cypher query to retrieve the entity with both Major and Minor kinds
	query := `
        MATCH (e {Id: $Id})
        RETURN labels(e)[0] AS MajorKind, e.MinorKind AS MinorKind, e.Id AS Id, e.Name AS Name, 
               toString(e.Created) AS Created, 
               CASE WHEN e.Terminated IS NOT NULL THEN toString(e.Terminated) ELSE NULL END AS Terminated
    `

	// Run the query
	result, err := session.Run(ctx, query, map[string]interface{}{"Id": entityID})
	if err != nil {
		log.Printf("[neo4j_client.ReadGraphEntity] error querying entity: %v", err)
		return nil, fmt.Errorf("error querying entity: %v", err)
	}

	// Process the result
	if result.Next(ctx) {
		record := result.Record()

		// Map the entity properties
		entity := map[string]interface{}{
			"Id":        fmt.Sprintf("%v", record.Values[2]), // e.Id
			"Name":      fmt.Sprintf("%v", record.Values[3]), // e.Name
			"Created":   fmt.Sprintf("%v", record.Values[4]), // e.Created
			"MajorKind": fmt.Sprintf("%v", record.Values[0]), // labels(e)[0]
			"MinorKind": fmt.Sprintf("%v", record.Values[1]), // e.MinorKind
		}

		// Add Terminated if it exists
		if terminatedVal, exists := record.Get("Terminated"); exists && terminatedVal != nil {
			entity["Terminated"] = fmt.Sprintf("%v", terminatedVal)
		}

		return entity, nil
	}

	// If no entity is found
	return nil, fmt.Errorf("entity with Id %s not found", entityID)
}

// ReadRelatedGraphEntityIds retrieves related entity IDs based on a given relationship
func (r *Neo4jRepository) ReadRelatedGraphEntityIds(ctx context.Context, entityID string, relationship string, ts string) ([]string, error) {
	if entityID == "" {
		return nil, fmt.Errorf("entity Id cannot be empty")
	}

	session := r.getSession(ctx)
	defer session.Close(ctx)

	query := fmt.Sprintf(`
        MATCH (e {Id: $entityID})-[r:%s]->(related)
        WHERE r.Created <= date($ts) AND (r.Terminated IS NULL OR r.Terminated > date($ts))
        RETURN related.Id AS relatedID
    `, relationship)

	result, err := session.Run(ctx, query, map[string]interface{}{
		"entityID": entityID,
		"ts":       ts,
	})
	if err != nil {
		log.Printf("[neo4j_client.ReadRelatedGraphEntityIds] error querying related entities: %v", err)
		return nil, fmt.Errorf("error querying related entities: %v", err)
	}

	var relatedIDs []string
	for result.Next(ctx) {
		record := result.Record()
		if relatedID, exists := record.Get("relatedID"); exists && relatedID != nil {
			relatedIDs = append(relatedIDs, fmt.Sprintf("%v", relatedID))
		}
	}

	if err := result.Err(); err != nil {
		log.Printf("[neo4j_client.ReadRelatedGraphEntityIds] error iterating over query result: %v", err)
		return nil, fmt.Errorf("error iterating over query result: %v", err)
	}

	return relatedIDs, nil
}

func (r *Neo4jRepository) ReadRelationships(ctx context.Context, entityID string) ([]map[string]interface{}, error) {

	if entityID == "" {
		return nil, fmt.Errorf("entity Id cannot be empty")
	}

	// Open session
	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Cypher query to get all relationships (incoming and outgoing)
	query := `
        MATCH (e {Id: $entityID})-[r]->(related)
        RETURN type(r) AS type, related.Id AS relatedID, "OUTGOING" AS direction, 
               toString(r.Created) AS Created, 
               CASE WHEN r.Terminated IS NOT NULL THEN toString(r.Terminated) ELSE NULL END AS Terminated,
               r.Id AS relationshipID
        UNION
        MATCH (e {Id: $entityID})<-[r]-(related)
        RETURN type(r) AS type, related.Id AS relatedID, "INCOMING" AS direction, 
               toString(r.Created) AS Created, 
               CASE WHEN r.Terminated IS NOT NULL THEN toString(r.Terminated) ELSE NULL END AS Terminated,
               r.Id AS relationshipID
    `

	// Run the query
	result, err := session.Run(ctx, query, map[string]interface{}{
		"entityID": entityID,
	})
	if err != nil {
		log.Printf("[neo4j_client.ReadRelationships] error querying relationships: %v", err)
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
			"Created":        fmt.Sprintf("%v", values[3]),
			"relationshipID": fmt.Sprintf("%v", values[5]), // Relationship ID
		}

		// Optional Terminated
		if values[4] != nil {
			rel["Terminated"] = fmt.Sprintf("%v", values[4])
		}

		relationships = append(relationships, rel)
	}

	// Return relationships as a map
	return relationships, nil
}

func (r *Neo4jRepository) ReadRelationship(ctx context.Context, relationshipID string) (map[string]interface{}, error) {

	if relationshipID == "" {
		return nil, fmt.Errorf("relationship Id cannot be empty")
	}

	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Cypher query to find the relationship by its ID
	query := `
        MATCH ()-[r]->()
        WHERE r.Id = $relationshipID
        RETURN type(r) AS type, startNode(r).Id AS startEntityID, endNode(r).Id AS endEntityID, 
               toString(r.Created) AS Created, 
               CASE WHEN r.Terminated IS NOT NULL THEN toString(r.Terminated) ELSE NULL END AS Terminated, 
               r.Id AS relationshipID
    `

	// Run the query to fetch the relationship
	result, err := session.Run(ctx, query, map[string]interface{}{
		"relationshipID": relationshipID,
	})
	if err != nil {
		log.Printf("[neo4j_client.ReadRelationship] error querying relationship: %v", err)
		return nil, fmt.Errorf("error querying relationship: %v", err)
	}

	// Process results
	if result.Next(ctx) {
		record := result.Record()
		values := record.Values

		// Ensure expected values exist
		if len(values) < 6 {
			log.Printf("[neo4j_client.ReadRelationship] unexpected data format for relationship")
			return nil, fmt.Errorf("unexpected data format for relationship")
		}

		// Map to hold the relationship data
		relationship := map[string]interface{}{
			"type":           fmt.Sprintf("%v", values[0]), // Relationship type
			"startEntityID":  fmt.Sprintf("%v", values[1]), // ID of the start entity
			"endEntityID":    fmt.Sprintf("%v", values[2]), // ID of the end entity
			"Created":        fmt.Sprintf("%v", values[3]),
			"relationshipID": fmt.Sprintf("%v", values[5]), // Relationship ID
		}

		// Optional Terminated
		if values[4] != nil {
			relationship["Terminated"] = fmt.Sprintf("%v", values[4])
		}

		// Return the relationship data as a map
		return relationship, nil
	}

	// If no relationship was found
	return nil, fmt.Errorf("relationship with Id %s not found", relationshipID)
}

// UpdateGraphEntity updates the properties of an existing entity
func (r *Neo4jRepository) UpdateGraphEntity(ctx context.Context, id string, updateData map[string]interface{}) (map[string]interface{}, error) {
	if id == "" {
		return nil, fmt.Errorf("entity Id cannot be empty")
	}

	// Prepare update parameters
	params := map[string]interface{}{
		"Id": id,
	}

	// Open session
	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Check if the entity exists
	existsQuery := `MATCH (e {Id: $Id}) RETURN e`
	result, err := session.Run(ctx, existsQuery, params)
	if err != nil {
		log.Printf("[neo4j_client.UpdateGraphEntity] error checking if entity exists: %v", err)
		return nil, fmt.Errorf("error checking if entity exists: %v", err)
	}

	if !result.Next(ctx) {
		log.Printf("[neo4j_client.UpdateGraphEntity] entity with Id %s does not exist", id)
		return nil, fmt.Errorf("entity with Id %s does not exist", id)
	}

	// Build Cypher query for updating entity
	query := `
        MATCH (e {Id: $Id})
    `

	// Add `Name` if provided
	if name, exists := updateData["Name"]; exists {
		params["Name"] = name
		query += `SET e.Name = $Name `
	}

	// Add `Terminated` if provided
	if terminated, exists := updateData["Terminated"]; exists {
		params["Terminated"] = terminated
		query += `SET e.Terminated = date($Terminated) `
	}

	// Execute update query and return updated entity
	query += ` RETURN e`

	result, err = session.Run(ctx, query, params)
	if err != nil {
		log.Printf("[neo4j_client.UpdateGraphEntity] error updating entity: %v", err)
		return nil, fmt.Errorf("error updating entity: %v", err)
	}

	// Retrieve updated entity
	if result.Next(ctx) {
		node, ok := result.Record().Get("e")
		if !ok {
			log.Printf("[neo4j_client.UpdateGraphEntity] unexpected error retrieving entity")
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

func (r *Neo4jRepository) UpdateRelationship(ctx context.Context, relationshipID string, updateData map[string]interface{}) (map[string]interface{}, error) {

	if relationshipID == "" {
		log.Printf("[neo4j_client.UpdateRelationship] relationship Id cannot be empty")
		return nil, fmt.Errorf("relationship Id cannot be empty")
	}

	// Prepare update parameters
	params := map[string]interface{}{
		"relationshipID": relationshipID,
	}

	// Open session
	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Check if the relationship exists
	existsQuery := `MATCH ()-[r {Id: $relationshipID}]->() RETURN r`
	result, err := session.Run(ctx, existsQuery, params)
	if err != nil {
		log.Printf("[neo4j_client.UpdateRelationship] error checking if relationship exists: %v", err)
		return nil, fmt.Errorf("error checking if relationship exists: %v", err)
	}

	if !result.Next(ctx) {
		log.Printf("[neo4j_client.UpdateRelationship] relationship with Id %s does not exist", relationshipID)
		return nil, fmt.Errorf("relationship with Id %s does not exist", relationshipID)
	}

	// Build Cypher query for updating relationship
	query := `
        MATCH ()-[r {Id: $relationshipID}]->()
    `

	// Add `Terminated` if provided (required)
	terminated, exists := updateData["Terminated"]
	if !exists {
		return nil, fmt.Errorf("terminated is required")
	}
	params["Terminated"] = terminated
	query += `SET r.Terminated = date($Terminated) RETURN r`

	// Execute update query and return updated relationship
	result, err = session.Run(ctx, query, params)
	if err != nil {
		log.Printf("[neo4j_client.UpdateRelationship] error updating relationship: %v", err)
		return nil, fmt.Errorf("error updating relationship: %v", err)
	}

	// Retrieve updated relationship
	if result.Next(ctx) {
		rel, ok := result.Record().Get("r")
		if !ok {
			log.Printf("[neo4j_client.UpdateRelationship] unexpected error retrieving relationship")
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

func (r *Neo4jRepository) DeleteRelationship(ctx context.Context, relationshipID string) error {
	if relationshipID == "" {
		return fmt.Errorf("entity Id cannot be empty")
	}

	// Prepare query parameters
	params := map[string]interface{}{
		"relationshipID": relationshipID,
	}

	// Open session
	session := r.getSession(ctx)
	defer session.Close(ctx)

	// Check if the relationship exists
	query := `MATCH ()-[r {Id: $relationshipID}]->() RETURN r`
	result, err := session.Run(ctx, query, params)
	if err != nil {
		log.Printf("[neo4j_client.DeleteRelationship] error checking if relationship exists: %v", err)
		return fmt.Errorf("error checking if relationship exists: %v", err)
	}

	// If no relationship is found, return an error
	if !result.Next(ctx) {
		log.Printf("[neo4j_client.DeleteRelationship] relationship with Id %s does not exist", relationshipID)
		return fmt.Errorf("relationship with Id %s does not exist", relationshipID)
	}

	// Delete the relationship
	deleteQuery := `MATCH ()-[r {Id: $relationshipID}]->() DELETE r`
	_, err = session.Run(ctx, deleteQuery, params)
	if err != nil {
		log.Printf("[neo4j_client.DeleteRelationship] error deleting relationship: %v", err)
		return fmt.Errorf("error deleting relationship: %v", err)
	}

	return nil
}

// DeleteGraphEntity deletes an entity by its ID
func (r *Neo4jRepository) DeleteGraphEntity(ctx context.Context, entityID string) error {
	if entityID == "" {
		log.Printf("[neo4j_client.DeleteGraphEntity] entity Id cannot be empty")
		return fmt.Errorf("entity Id cannot be empty")
	}

	session := r.getSession(ctx)
	defer session.Close(ctx)

	query := `MATCH (e {Id: $entityID}) RETURN e`
	params := map[string]interface{}{
		"entityID": entityID,
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		log.Printf("[neo4j_client.DeleteGraphEntity] error checking if entity exists: %v", err)
		return fmt.Errorf("error checking if entity exists: %v", err)
	}

	if !result.Next(ctx) {
		log.Printf("[neo4j_client.DeleteGraphEntity] entity with Id %s does not exist", entityID)
		return fmt.Errorf("entity with Id %s does not exist", entityID)
	}

	// Get the relationships of the entity
	relationships, err := r.ReadRelationships(ctx, entityID)
	if err != nil {
		log.Printf("[neo4j_client.DeleteGraphEntity] error getting relationships: %v", err)
		return fmt.Errorf("error getting relationships: %v", err)
	}

	// If there are relationships, return an error with relationship details
	if len(relationships) > 0 {
		log.Printf("[neo4j_client.DeleteGraphEntity] entity has relationships and cannot be deleted. Relationships: %v", relationships)
		return fmt.Errorf("entity has relationships and cannot be deleted. Relationships: %v", relationships)
	}

	// Delete the entity (node) with the given Id
	deleteQuery := `MATCH (e {Id: $entityID}) DELETE e`
	_, err = session.Run(ctx, deleteQuery, params)
	if err != nil {
		log.Printf("[neo4j_client.DeleteGraphEntity] error deleting entity: %v", err)
		return fmt.Errorf("error deleting entity: %v", err)
	}

	return nil
}
