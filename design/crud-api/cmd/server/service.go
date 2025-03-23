package main

import (
	"context"
	"log"
	"net"
	"os"

	"lk/datafoundation/crud-api/db/config"
	pb "lk/datafoundation/crud-api/lk/datafoundation/crud-api"

	"github.com/joho/godotenv"

	mongorepository "lk/datafoundation/crud-api/db/repository/mongo"
	neo4jrepository "lk/datafoundation/crud-api/db/repository/neo4j"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Server implements the CrudService
type Server struct {
	pb.UnimplementedCrudServiceServer
	mongoRepo *mongorepository.MongoRepository
	neo4jRepo *neo4jrepository.Neo4jRepository
}

// validateGraphEntityCreation checks if an entity has all required fields for Neo4j storage
func (s *Server) validateGraphEntityCreation(entity *pb.Entity) bool {
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

// CreateEntity handles entity creation with metadata
func (s *Server) CreateEntity(ctx context.Context, req *pb.Entity) (*pb.Entity, error) {
	log.Printf("Creating Entity: %s", req.Id)

	// Always save the entity in MongoDB, even if it has no metadata
	// The HandleMetadata function will only process it if it has metadata
	err := s.mongoRepo.HandleMetadata(ctx, req.Id, req)
	if err != nil {
		log.Printf("Error saving metadata in MongoDB: %v", err)
		return nil, err
	}
	log.Printf("Successfully saved metadata in MongoDB for entity: %s", req.Id)

	// Validate required fields for Neo4j entity creation
	if s.validateGraphEntityCreation(req) {
		// Prepare data for Neo4j with safety checks
		entityMap := map[string]interface{}{
			"Id": req.Id,
		}

		// Handle Kind field safely
		if req.Kind != nil {
			entityMap["Kind"] = req.Kind.GetMajor()
		}

		// Handle Name field safely
		if req.Name != nil && req.Name.GetValue() != nil {
			entityMap["Name"] = req.Name.GetValue().String()
		}

		// Handle other fields
		if req.Created != "" {
			entityMap["Created"] = req.Created
		}

		if req.Terminated != "" {
			entityMap["Terminated"] = req.Terminated
		}

		// Save entity in Neo4j only if validation passes
		_, err = s.neo4jRepo.CreateGraphEntity(ctx, entityMap)
		if err != nil {
			log.Printf("Error saving entity in Neo4j: %v", err)
			// Don't return error here so MongoDB storage is still successful
		} else {
			log.Printf("Successfully saved entity in Neo4j for entity: %s", req.Id)
		}
	} else {
		log.Printf("Entity %s saved in MongoDB only, skipping Neo4j due to missing required fields", req.Id)
	}

	return req, nil
}

// ReadEntity retrieves an entity's metadata
func (s *Server) ReadEntity(ctx context.Context, req *pb.Entity) (*pb.Entity, error) {
	log.Printf("Reading Entity metadata: %s", req.Id)
	debugMetadata(req)

	// Initialize metadata map
	var metadata map[string]*anypb.Any

	// Get the entity from MongoDB
	metadataResult, err := s.mongoRepo.GetMetadata(ctx, req.Id)
	if err != nil {
		// Log error and continue with empty metadata
		log.Printf("Error retrieving metadata for entity %s: %v", req.Id, err)
		metadata = make(map[string]*anypb.Any)
		log.Printf("METADATA DEBUG: Entity %s using empty metadata map: %+v", req.Id, metadata)
	} else {
		metadata = metadataResult
		log.Printf("METADATA DEBUG: Entity %s retrieved metadata: %+v", req.Id, metadata)
		// Print keys in metadata
		for k := range metadata {
			log.Printf("METADATA DEBUG: Entity %s has metadata key: %s", req.Id, k)
		}
	}

	// Try to get additional entity information from Neo4j
	var kind *pb.Kind
	var name *pb.TimeBasedValue
	var created string
	var terminated string

	// Attempt to read from Neo4j, but don't fail if it's not available
	entityMap, err := s.neo4jRepo.ReadGraphEntity(ctx, req.Id)
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

	// Return entity with all available fields
	return &pb.Entity{
		Id:            req.Id,
		Kind:          kind,
		Name:          name,
		Created:       created,
		Terminated:    terminated,
		Metadata:      metadata,
		Attributes:    make(map[string]*pb.TimeBasedValueList), // Empty attributes
		Relationships: make(map[string]*pb.Relationship),       // Empty relationships
	}, nil
}

// UpdateEntity modifies existing metadata
func (s *Server) UpdateEntity(ctx context.Context, req *pb.UpdateEntityRequest) (*pb.Entity, error) {
	// Extract ID from request parameter and entity data
	updateEntityID := req.Id
	updateEntity := req.Entity

	log.Printf("Updating Entity metadata: %s", updateEntityID)

	// Initialize metadata
	var metadata map[string]*anypb.Any

	// Pass the ID and metadata to HandleMetadata
	err := s.mongoRepo.HandleMetadata(ctx, updateEntityID, updateEntity)
	if err != nil {
		// Log error and continue with empty metadata
		log.Printf("Error updating metadata for entity %s: %v", updateEntityID, err)
		metadata = make(map[string]*anypb.Any)
	} else {
		// Use the provided metadata
		metadata = updateEntity.Metadata
	}

	// Return updated entity
	return &pb.Entity{
		Id:       updateEntity.Id,
		Metadata: metadata,
	}, nil
}

// DeleteEntity removes metadata
func (s *Server) DeleteEntity(ctx context.Context, req *pb.EntityId) (*pb.Empty, error) {
	log.Printf("Deleting Entity metadata: %s", req.Id)
	_, err := s.mongoRepo.DeleteEntity(ctx, req.Id)
	if err != nil {
		// Log error but return success
		log.Printf("Error deleting metadata for entity %s: %v", req.Id, err)
	}
	return &pb.Empty{}, nil
}

// Start the gRPC server
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize MongoDB config
	mongoConfig := &config.MongoConfig{
		URI:        os.Getenv("MONGO_URI"),
		DBName:     os.Getenv("MONGO_DB_NAME"),
		Collection: os.Getenv("MONGO_COLLECTION"),
	}

	// Initialize Neo4j config
	neo4jConfig := &config.Neo4jConfig{
		URI:      os.Getenv("NEO4J_URI"),
		Username: os.Getenv("NEO4J_USER"),
		Password: os.Getenv("NEO4J_PASSWORD"),
	}

	// Create MongoDB repository
	ctx := context.Background()
	mongoRepo := mongorepository.NewMongoRepository(ctx, mongoConfig)

	// Create Neo4j repository
	neo4jRepo, err := neo4jrepository.NewNeo4jRepository(ctx, neo4jConfig)
	if err != nil {
		log.Fatalf("Failed to create Neo4j repository: %v", err)
	}
	defer neo4jRepo.Close(ctx)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := &Server{
		mongoRepo: mongoRepo,
		neo4jRepo: neo4jRepo,
	}

	pb.RegisterCrudServiceServer(grpcServer, server)

	log.Println("CRUD Service is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
