package main

import (
	"context"
	"log"
	"net"
	"os"

	"lk/datafoundation/crud-api/db/config"
	pb "lk/datafoundation/crud-api/lk/datafoundation/crud-api"

	mongorepository "lk/datafoundation/crud-api/db/repository/mongo"
	neo4jrepository "lk/datafoundation/crud-api/db/repository/neo4j"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
)

// Server implements the CrudService
type Server struct {
	pb.UnimplementedCrudServiceServer
	mongoRepo *mongorepository.MongoRepository
	neo4jRepo *neo4jrepository.Neo4jRepository
}

// CreateEntity handles entity creation with metadata
func (s *Server) CreateEntity(ctx context.Context, req *pb.Entity) (*pb.Entity, error) {
	log.Printf("Creating Entity: %s", req.Id)

	// Always save the entity in MongoDB, even if it has no metadata
	// The HandleMetadata function will only process it if it has metadata
	err := s.mongoRepo.HandleMetadata(ctx, req.Id, req)
	if err != nil {
		log.Printf("[server.CreateEntity] Error saving metadata in MongoDB: %v", err)
		return nil, err
	} else {
		log.Printf("[server.CreateEntity] Successfully saved metadata in MongoDB for entity: %s", req.Id)
	}

	// Validate required fields for Neo4j entity creation
	success, err := s.neo4jRepo.HandleGraphEntityCreation(ctx, req)
	if !success {
		log.Printf("[server.CreateEntity] Error saving entity in Neo4j: %v", err)
		return nil, err
	} else {
		log.Printf("[server.CreateEntity] Successfully saved entity in Neo4j for entity: %s", req.Id)
	}

	// TODO: Add logic to handle relationships
	err = s.neo4jRepo.HandleGraphRelationshipsCreate(ctx, req)
	if err != nil {
		log.Printf("[server.CreateEntity] Error saving relationships in Neo4j: %v", err)
		return nil, err
	} else {
		log.Printf("[server.CreateEntity] Successfully saved relationships in Neo4j for entity: %s", req.Id)
	}

	// TODO: Add logic to handle attributes
	return req, nil
}

// ReadEntity retrieves an entity's metadata
func (s *Server) ReadEntity(ctx context.Context, req *pb.Entity) (*pb.Entity, error) {
	log.Printf("Reading Entity metadata: %s", req.Id)
	debugMetadata(req)

	// Get the entity from MongoDB - error handling and logging now in repository
	metadata, _ := s.mongoRepo.GetMetadata(ctx, req.Id)

	// Try to get additional entity information from Neo4j
	kind, name, created, terminated, _ := s.neo4jRepo.GetGraphEntity(ctx, req.Id)

	// Try to get relationships from Neo4j
	relationships, _ := s.neo4jRepo.GetGraphRelationships(ctx, req.Id)

	// Return entity with all available fields
	return &pb.Entity{
		Id:            req.Id,
		Kind:          kind,
		Name:          name,
		Created:       created,
		Terminated:    terminated,
		Metadata:      metadata,
		Attributes:    make(map[string]*pb.TimeBasedValueList), // Empty attributes
		Relationships: relationships,                           // Include relationships from Neo4j
	}, nil
}

// UpdateEntity modifies existing metadata
func (s *Server) UpdateEntity(ctx context.Context, req *pb.UpdateEntityRequest) (*pb.Entity, error) {
	// Extract ID from request parameter and entity data
	updateEntityID := req.Id
	updateEntity := req.Entity

	log.Printf("[server.UpdateEntity] Updating Entity: %s", updateEntityID)

	// Initialize metadata
	var metadata map[string]*anypb.Any

	// Pass the ID and metadata to HandleMetadata
	err := s.mongoRepo.HandleMetadata(ctx, updateEntityID, updateEntity)
	if err != nil {
		// Log error and continue with empty metadata
		log.Printf("[server.UpdateEntity] Error updating metadata for entity %s: %v", updateEntityID, err)
		metadata = make(map[string]*anypb.Any)
	} else {
		// Use the provided metadata
		metadata = updateEntity.Metadata
	}

	// Handle Graph Entity update if entity has required fields
	success, err := s.neo4jRepo.HandleGraphEntityUpdate(ctx, updateEntity)
	if !success {
		log.Printf("[server.UpdateEntity] Error updating graph entity for %s: %v", updateEntityID, err)
		// Continue processing despite error
	}

	// Handle Relationships update
	err = s.neo4jRepo.HandleGraphRelationshipsUpdate(ctx, updateEntity)
	if err != nil {
		log.Printf("[server.UpdateEntity] Error updating relationships for entity %s: %v", updateEntityID, err)
		// Continue processing despite error
	}

	// Read entity data from Neo4j to include in response
	kind, name, created, terminated, _ := s.neo4jRepo.GetGraphEntity(ctx, updateEntityID)

	// Get relationships from Neo4j
	relationships, _ := s.neo4jRepo.GetGraphRelationships(ctx, updateEntityID)

	// Return updated entity with all available information
	return &pb.Entity{
		Id:            updateEntity.Id,
		Kind:          kind,
		Name:          name,
		Created:       created,
		Terminated:    terminated,
		Metadata:      metadata,
		Attributes:    make(map[string]*pb.TimeBasedValueList), // Empty attributes
		Relationships: relationships,
	}, nil
}

// DeleteEntity removes metadata
func (s *Server) DeleteEntity(ctx context.Context, req *pb.EntityId) (*pb.Empty, error) {
	log.Printf("[server.DeleteEntity] Deleting Entity metadata: %s", req.Id)
	_, err := s.mongoRepo.DeleteEntity(ctx, req.Id)
	if err != nil {
		// Log error but return success
		log.Printf("[server.DeleteEntity] Error deleting metadata for entity %s: %v", req.Id, err)
	}
	// TODO: Implement Relationship Deletion in Neo4j
	// TODO: Implement Entity Deletion in Neo4j
	// TODO: Implement Attribute Deletion in Neo4j
	return &pb.Empty{}, nil
}

// Start the gRPC server
func main() {

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
		log.Fatalf("[service.main] Failed to create Neo4j repository: %v", err)
	}
	defer neo4jRepo.Close(ctx)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("[service.main] Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := &Server{
		mongoRepo: mongoRepo,
		neo4jRepo: neo4jRepo,
	}

	pb.RegisterCrudServiceServer(grpcServer, server)

	log.Println("[service.main] CRUD Service is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("[service.main] Failed to serve: %v", err)
	}
}
