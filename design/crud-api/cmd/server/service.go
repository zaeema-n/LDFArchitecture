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

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Server implements the CrudService
type Server struct {
	pb.UnimplementedCrudServiceServer
	repo *mongorepository.MongoRepository
}

// convertMetadata converts map[string]*anypb.Any to map[string]interface{}
func convertMetadata(in map[string]*anypb.Any) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range in {
		out[k] = v.GetValue()
	}
	return out
}

// CreateEntity handles entity creation with metadata
func (s *Server) CreateEntity(ctx context.Context, req *pb.Entity) (*pb.Entity, error) {
	log.Printf("Creating Entity with metadata: %s", req.Id)
	err := s.repo.HandleMetadata(ctx, req.Id, convertMetadata(req.Metadata))
	if err != nil {
		return nil, err
	}
	return req, nil
}

// ReadEntity retrieves an entity's metadata
func (s *Server) ReadEntity(ctx context.Context, req *pb.EntityId) (*pb.Entity, error) {
	log.Printf("Reading Entity metadata: %s", req.Id)
	metadata, err := s.repo.GetMetadata(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	// Convert back to Any
	anyMetadata := make(map[string]*anypb.Any)
	for k, v := range metadata {
		// Wrap string in a StringValue proto message
		anyVal, err := anypb.New(wrapperspb.String(v))
		if err != nil {
			return nil, err
		}
		anyMetadata[k] = anyVal
	}
	return &pb.Entity{
		Id:       req.Id,
		Metadata: anyMetadata,
	}, nil
}

// UpdateEntity modifies existing metadata
func (s *Server) UpdateEntity(ctx context.Context, req *pb.Entity) (*pb.Entity, error) {
	log.Printf("Updating Entity metadata: %s", req.Id)
	err := s.repo.HandleMetadata(ctx, req.Id, convertMetadata(req.Metadata))
	if err != nil {
		return nil, err
	}
	return req, nil
}

// DeleteEntity removes metadata
func (s *Server) DeleteEntity(ctx context.Context, req *pb.EntityId) (*pb.Empty, error) {
	log.Printf("Deleting Entity metadata: %s", req.Id)
	_, err := s.repo.DeleteEntity(ctx, req.Id)
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
	server := &Server{repo: repo}

	pb.RegisterCrudServiceServer(grpcServer, server)

	log.Println("CRUD Service is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
