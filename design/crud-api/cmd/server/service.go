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
)

// Server implements the CrudService
type Server struct {
	pb.UnimplementedCrudServiceServer
	mongoRepo *mongorepository.MongoRepository
	neo4jRepo *neo4jrepository.Neo4jRepository
}

// type Neo4jServer struct {
// 	pb.UnimplementedCrudServiceServer
// 	repo *neo4jrepository.Neo4jRepository
// }

// CreateEntity handles entity creation with metadata
func (s *Server) CreateEntity(ctx context.Context, req *pb.Entity) (*pb.Entity, error) {
	log.Printf("Creating Entity with metadata: %s", req.Id)
	err := s.mongoRepo.HandleMetadata(ctx, req.Id, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// ReadEntity retrieves an entity's metadata
func (s *Server) ReadEntity(ctx context.Context, req *pb.EntityId) (*pb.Entity, error) {
	log.Printf("Reading Entity metadata: %s", req.Id)
	metadata, err := s.mongoRepo.GetMetadata(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	// Convert back to Any

	return &pb.Entity{
		Id:       req.Id,
		Metadata: metadata,
	}, nil
}

// UpdateEntity modifies existing metadata
func (s *Server) UpdateEntity(ctx context.Context, req *pb.UpdateEntityRequest) (*pb.Entity, error) {
	// Extract ID from request parameter and entity data
	updateEntityID := req.Id
	updateEntity := req.Entity

	log.Printf("Updating Entity metadata: %s", updateEntityID)

	// Pass the ID and metadata to HandleMetadata
	err := s.mongoRepo.HandleMetadata(ctx, updateEntityID, updateEntity)
	if err != nil {
		return nil, err
	}

	// Return updated entity
	return &pb.Entity{
		Id:       updateEntity.Id,
		Metadata: updateEntity.Metadata,
	}, nil
}

// DeleteEntity removes metadata
func (s *Server) DeleteEntity(ctx context.Context, req *pb.EntityId) (*pb.Empty, error) {
	log.Printf("Deleting Entity metadata: %s", req.Id)
	_, err := s.mongoRepo.DeleteEntity(ctx, req.Id)
	if err != nil {
		return nil, err
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

	// Create MongoDB repository
	ctx := context.Background()
	repo := mongorepository.NewMongoRepository(ctx, mongoConfig)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := &Server{mongoRepo: repo}

	pb.RegisterCrudServiceServer(grpcServer, server)

	log.Println("CRUD Service is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
