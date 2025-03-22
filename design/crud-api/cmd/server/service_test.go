package main

import (
	"context"
	"log"
	"os"
	"testing"

	"lk/datafoundation/crud-api/db/config"
	mongorepository "lk/datafoundation/crud-api/db/repository/mongo"
	neo4jrepository "lk/datafoundation/crud-api/db/repository/neo4j"
)

var server *Server

// TestMain sets up the actual MongoDB and Neo4j repositories before running the tests
func TestMain(m *testing.M) {
	// Load environment variables for database configurations
	neo4jConfig := &config.Neo4jConfig{
		URI:      os.Getenv("NEO4J_URI"),
		Username: os.Getenv("NEO4J_USER"),
		Password: os.Getenv("NEO4J_PASSWORD"),
	}

	mongoConfig := &config.MongoConfig{
		URI:        os.Getenv("MONGO_URI"),
		DBName:     os.Getenv("MONGO_DB_NAME"),
		Collection: os.Getenv("MONGO_COLLECTION"),
	}

	// Initialize Neo4j repository
	ctx := context.Background()
	neo4jRepo, err := neo4jrepository.NewNeo4jRepository(ctx, neo4jConfig)
	if err != nil {
		log.Fatalf("Failed to initialize Neo4j repository: %v", err)
	}
	defer neo4jRepo.Close(ctx)

	// Initialize MongoDB repository
	mongoRepo := mongorepository.NewMongoRepository(ctx, mongoConfig)
	if mongoRepo == nil {
		log.Fatalf("Failed to initialize MongoDB repository")
	}

	// Create the server with the initialized repositories
	server = &Server{
		mongoRepo: mongoRepo,
		neo4jRepo: neo4jRepo,
	}

	// Run the tests
	code := m.Run()

	// Exit with the test result code
	os.Exit(code)
}

// TestCreateEntity tests the CreateEntity function
// TODO: FIX BUG @zaeema-n there is a bug in the CreateEntity function
func TestCreateEntity(t *testing.T) {

	// Encode Name as *anypb.Any
	// nameValue, err := anypb.New(&anypb.Any{Value: []byte("John Doe")})
	// assert.NoError(t, err)

	// // Define test cases
	// tests := []struct {
	// 	name    string
	// 	entity  *pb.Entity
	// 	wantErr bool
	// }{
	// 	{
	// 		name: "CreateEntity_Success",
	// 		entity: &pb.Entity{
	// 			Id: "123",
	// 			Kind: &pb.Kind{
	// 				Major: "Person",
	// 				Minor: "Employee",
	// 			},
	// 			Name:       &pb.TimeBasedValue{Value: nameValue},
	// 			Created:    "2025-03-20",
	// 			Terminated: "",
	// 		},
	// 		wantErr: false,
	// 	},
	// 	{
	// 		name: "CreateEntity_MissingId",
	// 		entity: &pb.Entity{
	// 			Id: "",
	// 			Kind: &pb.Kind{
	// 				Major: "Person",
	// 				Minor: "Employee",
	// 			},
	// 			Name:       &pb.TimeBasedValue{Value: nameValue},
	// 			Created:    "2025-03-20",
	// 			Terminated: "",
	// 		},
	// 		wantErr: true,
	// 	},
	// }

	// // Run test cases
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		ctx := context.Background()

	// 		resp, err := server.CreateEntity(ctx, tt.entity)
	// 		if (err != nil) != tt.wantErr {
	// 			t.Errorf("CreateEntity() error = %v, wantErr %v", err, tt.wantErr)
	// 		}

	// 		if !tt.wantErr {
	// 			// Verify the response matches the input entity
	// 			assert.Equal(t, tt.entity, resp, "Expected response to match the input entity")
	// 		}
	// 	})
	// }
}
