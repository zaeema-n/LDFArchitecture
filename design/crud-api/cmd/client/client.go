package main

import (
	"context"
	"log"
	"time"

	pb "lk/datafoundation/crud-api/lk/datafoundation/crud-api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Define a timeout for establishing the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a new gRPC client using NewClient
	clientConn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()), // Insecure for local development
	)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer clientConn.Close()

	client := pb.NewCrudServiceClient(clientConn)

	// Create a new entity
	entity := &pb.Entity{
		Id: "123",
		Kind: &pb.Kind{
			Major: "example",
			Minor: "test",
		},
		Created: "2023-01-01T00:00:00Z",
	}

	// Call CreateEntity
	createdEntity, err := client.CreateEntity(ctx, entity)
	if err != nil {
		log.Fatalf("Could not create entity: %v", err)
	}
	log.Printf("Created Entity: %v", createdEntity)

	// Call ReadEntity
	readEntity, err := client.ReadEntity(ctx, &pb.Entity{Id: "123"})
	if err != nil {
		log.Fatalf("Could not read entity: %v", err)
	}
	log.Printf("Read Entity: %v", readEntity)

	// Call UpdateEntity
	entity.Kind.Minor = "updated"
	updatedEntity, err := client.UpdateEntity(ctx, entity)
	if err != nil {
		log.Fatalf("Could not update entity: %v", err)
	}
	log.Printf("Updated Entity: %v", updatedEntity)

	// Call DeleteEntity
	deletedEntity, err := client.DeleteEntity(ctx, &pb.Entity{Id: "123"})
	if err != nil {
		log.Fatalf("Could not delete entity: %v", err)
	}
	log.Printf("Deleted Entity: %v", deletedEntity)
}
