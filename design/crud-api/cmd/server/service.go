package main

import (
	"context"
	"log"
	"net"

	pb "lk/datafoundation/crud-api/lk/datafoundation/crud-api"

	"google.golang.org/grpc"
)

// Server implements the CrudService
type Server struct {
	pb.UnimplementedCrudServiceServer
	entities map[string]*pb.Entity
}

// CreateEntity handles entity creation
func (s *Server) CreateEntity(ctx context.Context, req *pb.Entity) (*pb.Entity, error) {
	log.Printf("Creating Entity: %s", req.Id)
	s.entities[req.Id] = req
	return req, nil
}

// ReadEntity retrieves an entity
func (s *Server) ReadEntity(ctx context.Context, req *pb.Entity) (*pb.Entity, error) {
	log.Printf("Reading Entity: %s", req.Id)
	entity, exists := s.entities[req.Id]
	if !exists {
		return nil, nil
	}
	return entity, nil
}

// UpdateEntity modifies an existing entity
func (s *Server) UpdateEntity(ctx context.Context, req *pb.Entity) (*pb.Entity, error) {
	log.Printf("Updating Entity: %s", req.Id)
	s.entities[req.Id] = req
	return req, nil
}

// DeleteEntity removes an entity
func (s *Server) DeleteEntity(ctx context.Context, req *pb.Entity) (*pb.Entity, error) {
	log.Printf("Deleting Entity: %s", req.Id)
	delete(s.entities, req.Id)
	return req, nil
}

// Start the gRPC server
func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := &Server{entities: make(map[string]*pb.Entity)}

	pb.RegisterCrudServiceServer(grpcServer, server)

	log.Println("CRUD Service is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
