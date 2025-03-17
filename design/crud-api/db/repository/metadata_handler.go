package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Add this function to handle metadata operations
func (repo *MongoRepository) HandleMetadata(ctx context.Context, entityId string, metadata map[string]string) error {
	update := bson.M{
		"$set": bson.M{
			"metadata": metadata,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := repo.collection().UpdateOne(
		ctx,
		bson.M{"_id": entityId},
		update,
		opts,
	)
	return err
}

// Add a function to retrieve metadata
func (repo *MongoRepository) GetMetadata(ctx context.Context, entityId string) (map[string]string, error) {
	var result struct {
		Metadata map[string]string `bson:"metadata"`
	}

	err := repo.collection().FindOne(
		ctx,
		bson.M{"_id": entityId},
	).Decode(&result)

	return result.Metadata, err
}
