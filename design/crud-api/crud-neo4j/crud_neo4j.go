package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Entity struct defines the structure for the incoming JSON entity
type Entity struct {
	ID          string  `json:"id"`
	Kind        string  `json:"kind"`
	Name        string  `json:"name"`
	DateCreated string  `json:"dateCreated"`
	DateEnded   *string `json:"dateEnded,omitempty"`
}

// Relationship struct for JSON input
type Relationship struct {
	RelationshipID string  `json:"relationshipID"`
	ParentID       string  `json:"parentID"`
	ChildID        string  `json:"childID"`
	RelType        string  `json:"relType"`
	StartDate      string  `json:"startDate"`
	EndDate        *string `json:"endDate"` // Pointer for optional values
}

// Neo4jClient struct to manage the database connection
type Neo4jClient struct {
	Driver neo4j.Driver
}

// NewNeo4jClient initializes a new Neo4j client
func NewNeo4jClient(uri, username, password string) (*Neo4jClient, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("error creating Neo4j driver: %v", err)
	}
	return &Neo4jClient{Driver: driver}, nil
}

// Close the Neo4j connection
func (c *Neo4jClient) Close() {
	c.Driver.Close()
}

// getSession creates a new session
func (c *Neo4jClient) getSession() neo4j.Session {
	return c.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
}

// createEntity checks if an entity exists and creates it if it doesn’t
func (c *Neo4jClient) createEntity(entityJSON string) error {
	// Unmarshal the JSON string into an Entity struct
	var entity Entity
	err := json.Unmarshal([]byte(entityJSON), &entity)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	// Convert ID to integer
	entityID, err := strconv.Atoi(entity.ID)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	session := c.getSession()
	defer session.Close()

	// Check if the node already exists
	existsQuery := `MATCH (e:` + entity.Kind + ` {id: $id}) RETURN e`
	result, err := session.Run(existsQuery, map[string]interface{}{"id": entity.ID})
	if err != nil {
		return fmt.Errorf("error checking if entity exists: %v", err)
	}

	// If entity exists, return an error
	if result.Next() {
		return fmt.Errorf("entity with ID %s already exists", entity.ID)
	}

	// Create the node
	createQuery := `CREATE (e:` + entity.Kind + ` {id: $id, name: $name, dateCreated: date($dateCreated)`
	if entity.DateEnded != nil {
		createQuery += `, dateEnded: date($dateEnded)`
	}
	createQuery += `})`

	params := map[string]interface{}{
		"id":          entityID,
		"name":        entity.Name,
		"dateCreated": entity.DateCreated,
	}
	if entity.DateEnded != nil {
		params["dateEnded"] = *entity.DateEnded
	}

	_, err = session.Run(createQuery, params)
	if err != nil {
		return fmt.Errorf("error creating entity: %v", err)
	}

	return nil
}

func (c *Neo4jClient) createRelationship(relJSON string) error {
	// Unmarshal JSON into Relationship struct
	var rel Relationship
	err := json.Unmarshal([]byte(relJSON), &rel)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	// Convert IDs to integers
	parentInt, err := strconv.Atoi(rel.ParentID)
	if err != nil {
		return fmt.Errorf("invalid parent ID format: %v", err)
	}

	childInt, err := strconv.Atoi(rel.ChildID)
	if err != nil {
		return fmt.Errorf("invalid child ID format: %v", err)
	}

	relationshipInt, err := strconv.Atoi(rel.RelationshipID)
	if err != nil {
		return fmt.Errorf("invalid relationship ID format: %v", err)
	}

	session := c.getSession()
	defer session.Close()

	// Check if both entities exist
	existsQuery := `MATCH (p {id: $parentID}), (c {id: $childID}) RETURN p, c`
	result, err := session.Run(existsQuery, map[string]interface{}{
		"parentID": parentInt,
		"childID":  childInt,
	})
	if err != nil {
		return fmt.Errorf("error checking entities: %v", err)
	}
	if !result.Next() {
		return fmt.Errorf("either parent or child entity does not exist")
	}

	// Create the relationship
	createQuery := `MATCH (p {id: $parentID}), (c {id: $childID})
                    MERGE (p)-[r:` + rel.RelType + ` {id: $relationshipID}]->(c)
                    SET r.startDate = date($startDate)`
	if rel.EndDate != nil {
		createQuery += `, r.endDate = date($endDate)`
	}

	params := map[string]interface{}{
		"parentID":       parentInt,
		"childID":        childInt,
		"relationshipID": relationshipInt,
		"startDate":      rel.StartDate, // String, converted in Cypher
	}
	if rel.EndDate != nil {
		params["endDate"] = *rel.EndDate
	}

	_, err = session.Run(createQuery, params)
	if err != nil {
		return fmt.Errorf("error creating relationship: %v", err)
	}

	return nil
}

func (c *Neo4jClient) readEntity(entityID string) (string, error) {
	// Convert entity ID to an integer
	entityInt, err := strconv.Atoi(entityID)
	if err != nil {
		return "", fmt.Errorf("invalid entity ID format: %v", err)
	}

	session := c.getSession()
	defer session.Close()

	// Query to retrieve the entity details, including its label (kind)
	query := `MATCH (e {id: $id})
			  RETURN labels(e)[0] AS kind, e.id AS id, e.name AS name, 
			         toString(e.dateCreated) AS dateCreated, 
			         CASE WHEN e.dateEnded IS NOT NULL THEN toString(e.dateEnded) ELSE NULL END AS dateEnded`

	// Run the query
	result, err := session.Run(query, map[string]interface{}{"id": entityInt})
	if err != nil {
		return "", fmt.Errorf("error querying entity: %v", err)
	}

	// Extract entity data
	if result.Next() {
		record := result.Record()
		values := record.Values // Access the slice of values

		// Ensure values exist before accessing indexes
		if len(values) < 5 {
			return "", fmt.Errorf("unexpected query result format")
		}

		entity := Entity{
			ID:          fmt.Sprintf("%v", values[1]), // Convert ID back to string
			Kind:        fmt.Sprintf("%v", values[0]),
			Name:        fmt.Sprintf("%v", values[2]),
			DateCreated: fmt.Sprintf("%v", values[3]),
		}

		// Handle optional dateEnded field
		if values[4] != nil {
			dateEnded := fmt.Sprintf("%v", values[4])
			entity.DateEnded = &dateEnded
		}

		// Convert to JSON string
		jsonData, err := json.Marshal(entity)
		if err != nil {
			return "", fmt.Errorf("error marshalling entity to JSON: %v", err)
		}

		return string(jsonData), nil
	}

	return "", fmt.Errorf("entity with ID %s not found", entityID)
}

// getRelatedEntityIds retrieves related entity IDs based on a given relationship
func (c *Neo4jClient) getRelatedEntityIds(entityID, relationship string, ts string) ([]string, error) {
	// Convert entity ID to integer
	entityInt, err := strconv.Atoi(entityID)
	if err != nil {
		return nil, fmt.Errorf("invalid entity ID format: %v", err)
	}

	// Open session
	session := c.getSession()
	defer session.Close()

	// Query to find related entities within the given timestamp
	query := fmt.Sprintf(`
		MATCH (e {id: $entityID})-[r:%s]->(related)
		WHERE r.startDate <= date($ts) AND (r.endDate IS NULL OR r.endDate > date($ts))
		RETURN related.id AS relatedID
	`, relationship) // Using dynamic relationship type

	// Run the query
	result, err := session.Run(query, map[string]interface{}{
		"entityID": entityInt,
		"ts":       ts, // Date format should be "YYYY-MM-DD"
	})
	if err != nil {
		return nil, fmt.Errorf("error querying related entities: %v", err)
	}

	// Extract IDs
	var relatedIDs []string
	for result.Next() {
		record := result.Record()
		id, ok := record.Get("relatedID")
		if !ok {
			continue
		}
		relatedIDs = append(relatedIDs, fmt.Sprintf("%v", id))
	}

	return relatedIDs, nil
}

// getRelationships retrieves all relationships, incoming or outgoing, for a given entity
func (c *Neo4jClient) getRelationships(entityID string) (string, error) {
	// Convert entity ID to integer
	entityInt, err := strconv.Atoi(entityID)
	if err != nil {
		return "", fmt.Errorf("invalid entity ID format: %v", err)
	}

	// Open session
	session := c.getSession()
	defer session.Close()

	// Cypher query to get all relationships with direction
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

	// Run query
	result, err := session.Run(query, map[string]interface{}{
		"entityID": entityInt,
	})
	if err != nil {
		return "", fmt.Errorf("error querying relationships: %v", err)
	}

	// Process results
	var relationships []map[string]interface{}
	for result.Next() {
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

	// Convert relationships to JSON
	jsonData, err := json.MarshalIndent(relationships, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling relationships to JSON: %v", err)
	}

	return string(jsonData), nil
}

// can only update name and dateEnded
func (c *Neo4jClient) updateEntity(updateJSON string) error {
	// Parse JSON input
	var updateData map[string]interface{}
	err := json.Unmarshal([]byte(updateJSON), &updateData)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	// Extract entity ID
	entityID, ok := updateData["id"].(string)
	if !ok || entityID == "" {
		return fmt.Errorf("invalid or missing entity ID")
	}

	// Convert entity ID to integer
	entityInt, err := strconv.Atoi(entityID)
	if err != nil {
		return fmt.Errorf("invalid entity ID format: %v", err)
	}

	// Prepare update parameters
	params := map[string]interface{}{
		"id": entityInt,
	}

	// Open session
	session := c.getSession()
	defer session.Close()

	// Check if the entity exists
	existsQuery := `MATCH (e {id: $id}) RETURN e`
	result, err := session.Run(existsQuery, params)
	if err != nil {
		return fmt.Errorf("error checking if entity exists: %v", err)
	}

	if !result.Next() {
		return fmt.Errorf("entity with ID %d does not exist", entityInt)
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

	// Execute query
	_, err = session.Run(query, params)
	if err != nil {
		return fmt.Errorf("error updating entity: %v", err)
	}

	return nil
}

func (c *Neo4jClient) updateRelationship(updateJSON string) error {
	// Parse JSON input
	var updateData map[string]interface{}
	err := json.Unmarshal([]byte(updateJSON), &updateData)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	// Extract relationship ID
	relationshipID, ok := updateData["relationshipID"].(string)
	if !ok || relationshipID == "" {
		return fmt.Errorf("invalid or missing relationship ID")
	}

	// Convert relationship ID to integer
	relationshipInt, err := strconv.Atoi(relationshipID)
	if err != nil {
		return fmt.Errorf("invalid relationship ID format: %v", err)
	}

	// Prepare update parameters
	params := map[string]interface{}{
		"relationshipID": relationshipInt,
	}

	// Open session
	session := c.getSession()
	defer session.Close()

	// Check if the relationship exists
	existsQuery := `MATCH ()-[r {id: $relationshipID}]->() RETURN r`
	result, err := session.Run(existsQuery, params)
	if err != nil {
		return fmt.Errorf("error checking if relationship exists: %v", err)
	}

	if !result.Next() {
		return fmt.Errorf("relationship with ID %d does not exist", relationshipInt)
	}

	// Build Cypher query for updating relationship
	query := `
        MATCH ()-[r {id: $relationshipID}]->()
    `

	// Add `endDate` if provided
	if endDate, exists := updateData["endDate"]; exists {
		params["endDate"] = endDate
		query += `SET r.endDate = date($endDate) `
	} else {
		return fmt.Errorf("endDate is required")
	}

	// Execute query
	_, err = session.Run(query, params)
	if err != nil {
		return fmt.Errorf("error updating relationship: %v", err)
	}

	return nil
}

func (c *Neo4jClient) deleteRelationship(relationshipID string) error {
	// Convert relationship ID to integer
	relationshipInt, err := strconv.Atoi(relationshipID)
	if err != nil {
		return fmt.Errorf("invalid relationship ID format: %v", err)
	}

	// Open a session
	session := c.getSession()
	defer session.Close()

	// Check if the relationship exists
	query := `MATCH ()-[r {id: $relationshipID}]->() RETURN r`
	params := map[string]interface{}{
		"relationshipID": relationshipInt,
	}

	// Execute query to find the relationship
	result, err := session.Run(query, params)
	if err != nil {
		return fmt.Errorf("error checking if relationship exists: %v", err)
	}

	// If no relationship is found, return an error
	if !result.Next() {
		return fmt.Errorf("relationship with ID %d does not exist", relationshipInt)
	}

	// Delete the relationship
	deleteQuery := `MATCH ()-[r {id: $relationshipID}]->() DELETE r`
	_, err = session.Run(deleteQuery, params)
	if err != nil {
		return fmt.Errorf("error deleting relationship: %v", err)
	}

	return nil
}

// This only works if the entity has no relationships
func (c *Neo4jClient) deleteEntity(entityID string) error {
	// Convert entity ID to integer
	entityInt, err := strconv.Atoi(entityID)
	if err != nil {
		return fmt.Errorf("invalid entity ID format: %v", err)
	}

	// Open a session
	session := c.getSession()
	defer session.Close()

	// Check if the node exists before attempting to delete
	query := `MATCH (e {id: $entityID}) RETURN e`
	params := map[string]interface{}{
		"entityID": entityInt,
	}

	// Run query to find the entity
	result, err := session.Run(query, params)
	if err != nil {
		return fmt.Errorf("error checking if entity exists: %v", err)
	}

	// If entity doesn't exist, return an error
	if !result.Next() {
		return fmt.Errorf("entity with ID %d does not exist", entityInt)
	}

	// Get the relationships of the entity
	relationships, err := c.getRelationships(entityID)
	if err != nil {
		return fmt.Errorf("error getting relationships: %v", err)
	}

	// If there are relationships, return an error with relationship details
	if relationships != "null" { // Assuming relationships are returned as JSON array
		return fmt.Errorf("entity has relationships and cannot be deleted. Relationships: %s", relationships)
	}

	// Delete the entity (node) with the given ID
	deleteQuery := `MATCH (e {id: $entityID}) DELETE e`
	_, err = session.Run(deleteQuery, params)
	if err != nil {
		return fmt.Errorf("error deleting entity: %v", err)
	}

	return nil
}

func main() {

	// Load the environment variables from the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	neo4jURI := os.Getenv("NEO4J_URI")
	neo4jUsername := os.Getenv("NEO4J_USER")
	neo4jPassword := os.Getenv("NEO4J_PASSWORD")

	// Initialize Neo4j client
	client, err := NewNeo4jClient(neo4jURI, neo4jUsername, neo4jPassword)
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer client.Close()

	// // Example JSON entity with dateEnded
	// entityJSONWithDateEnded := `{
	//     "id": "1234",
	//     "kind": "person",
	//     "name": "John Doe",
	//     "dateCreated": "2025-03-17",
	//     "dateEnded": "2025-03-18"
	// }`

	// // Example JSON entity without dateEnded
	// entityJSONWithoutDateEnded := `{
	//     "id": "1235",
	//     "kind": "person",
	//     "name": "Jane Doe",
	//     "dateCreated": "2025-03-17"
	// }`

	// // Call createEntity function with dateEnded
	// err = client.createEntity(entityJSONWithDateEnded)
	// if err != nil {
	// 	log.Fatalf("Error creating entity: %v", err)
	// }
	// fmt.Println("Entity with dateEnded created successfully.")

	// // Call createEntity function without dateEnded
	// err = client.createEntity(entityJSONWithoutDateEnded)
	// if err != nil {
	// 	log.Fatalf("Error creating entity: %v", err)
	// }
	// fmt.Println("Entity without dateEnded created successfully.")

	// relInput := `{
	// 	"relationshipID": "1234",
	// 	"parentID": "1234",
	// 	"childID": "1235",
	// 	"relType": "PARENT_OF",
	// 	"startDate": "2025-03-17",
	// 	"endDate": "2025-12-31"
	// }`

	// relInput := `{
	// 	"relationshipID": "1235",
	// 	"parentID": "1235",
	// 	"childID": "1234",
	// 	"relType": "PARENT_OF",
	// 	"startDate": "2025-03-17"
	// }`

	// err = client.createRelationship(relInput)
	// if err != nil {
	// 	log.Println("Error:", err)
	// } else {
	// 	log.Println("Relationship created successfully!")
	// }

	// entityID := "1235"
	// entityJSON, err := client.readEntity(entityID)
	// if err != nil {
	// 	log.Fatalf("Error reading entity: %v", err)
	// }
	// fmt.Printf("Entity with ID %s: %s\n", entityID, entityJSON)

	// Test getRelatedEntityIds function
	// relatedEntityID := "1235"
	// relationshipType := "PARENT_OF"
	// timestamp := "2026-03-17"
	// relatedIDs, err := client.getRelatedEntityIds(relatedEntityID, relationshipType, timestamp)
	// if err != nil {
	// 	log.Fatalf("Error getting related entity IDs: %v", err)
	// }
	// fmt.Printf("Related entity IDs for entity %s with relationship %s at %s: %v\n", relatedEntityID, relationshipType, timestamp, relatedIDs)

	// entityID := "1234"
	// relationshipsJSON, err := client.getRelationships(entityID)
	// if err != nil {
	// 	log.Fatalf("Error getting relationships: %v", err)
	// }
	// fmt.Printf("Relationships for entity with ID %s: %s\n", entityID, relationshipsJSON)

	// updateJSON := `{
	// 	"id": "1245",
	// 	"dateEnded": "2026-03-18"
	// }`

	// err = client.updateEntity(updateJSON)
	// if err != nil {
	// 	log.Fatalf("Error updating entity: %v", err)
	// }
	// fmt.Println("Entity updated successfully.")

	// Test updateRelationship function
	// updateRelJSON := `{
	//     "relationshipID": "1245",
	//     "endDate": "2026-03-18"
	// }`

	// err = client.updateRelationship(updateRelJSON)
	// if err != nil {
	// 	log.Fatalf("Error updating relationship: %v", err)
	// }
	// fmt.Println("Relationship updated successfully.")

	// err = client.deleteRelationship("1235")
	// if err != nil {
	// 	log.Fatalf("Error deleting relationship: %v", err)
	// }
	// fmt.Println("Relationship deleted successfully.")

	// err = client.deleteEntity("1234")
	// if err != nil {
	// 	log.Fatalf("Error deleting entity: %v", err)
	// }
	// fmt.Println("Entity deleted successfully.")

}
