package repository

import (
	"context"

	pb "lk/datafoundation/crud-api/lk/datafoundation/crud-api"

	"google.golang.org/protobuf/types/known/anypb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Add this function to handle metadata operations
func (repo *MongoRepository) HandleMetadata(ctx context.Context, entityId string, metadata map[string]*anypb.Any) error {
	// Check if entity exists
	entity, err := repo.ReadEntity(ctx, entityId)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if entity == nil {
		// Create new entity with metadata
		newEntity := &pb.Entity{
			Id:       entityId,
			Metadata: metadata,
		}
		_, err = repo.CreateEntity(ctx, newEntity)
	} else {
		// Update existing entity's metadata
		_, err = repo.UpdateEntity(ctx, entityId, bson.M{"metadata": metadata})
	}

	return err
}

// Improved GetMetadata function that handles conversion internally
func (repo *MongoRepository) GetMetadata(ctx context.Context, entityId string) (map[string]*anypb.Any, error) {
	// Use the existing ReadEntity method for consistency
	entity, err := repo.ReadEntity(ctx, entityId)
	if err != nil {
		return nil, err
	}

	// Return the original protobuf Any metadata
	return entity.Metadata, nil
}
