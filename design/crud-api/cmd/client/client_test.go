package main

import (
	"context"
	"net"
	"testing"
	"time"

	pb "lk/datafoundation/crud-api/lk/datafoundation/crud-api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Mock server implementation
type mockServer struct {
	pb.UnimplementedCrudServiceServer
}

func (s *mockServer) CreateEntity(ctx context.Context, req *pb.Entity) (*pb.Entity, error) {
	return req, nil
}

func TestClient(t *testing.T) {
	// Start a mock gRPC server
	lis, err := net.Listen("tcp", ":0") // Use :0 to get a free port
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCrudServiceServer(grpcServer, &mockServer{})
	go grpcServer.Serve(lis)
	defer grpcServer.Stop()

	// Create a client connection to the mock server using grpc.DialContext
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewCrudServiceClient(conn)

	// Define test cases
	tests := []struct {
		name    string
		entity  *pb.Entity
		wantErr bool
	}{
		{
			name: "CreateEntity_Success",
			entity: &pb.Entity{
				Id: "123",
				Kind: &pb.Kind{
					Major: "example",
					Minor: "test",
				},
				Created: "2023-01-01T00:00:00Z",
			},
			wantErr: false,
		},
		// Add more test cases as needed
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			_, err := client.CreateEntity(ctx, tt.entity)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEntity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}