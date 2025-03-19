package mongorepository

import (
	"context"
	"lk/datafoundation/crud-api/db/config"
	"log"

	pb "lk/datafoundation/crud-api/lk/datafoundation/crud-api"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	client *mongo.Client
	config *config.MongoConfig
}

// NewMongoRepository initializes a MongoDB client
func NewMongoRepository(ctx context.Context, config *config.MongoConfig) *MongoRepository {
	clientOptions := options.Client().ApplyURI(config.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return &MongoRepository{
		client: client,
		config: config,
	}
}

func (repo *MongoRepository) collection() *mongo.Collection {
	return repo.client.Database(repo.config.DBName).Collection(repo.config.Collection)
}

// CreateEntity inserts a new entity in MongoDB
func (repo *MongoRepository) CreateEntity(ctx context.Context, entity *pb.Entity) (*mongo.InsertOneResult, error) {
	result, err := repo.collection().InsertOne(ctx, entity)
	return result, err
}

// ReadEntity fetches an entity by ID from MongoDB
func (repo *MongoRepository) ReadEntity(ctx context.Context, id string) (*pb.Entity, error) {
	var entity *pb.Entity
	err := repo.collection().FindOne(ctx, bson.M{"id": id}).Decode(&entity)
	return entity, err
}

// UpdateEntity updates an entity's attributes in MongoDB
func (repo *MongoRepository) UpdateEntity(ctx context.Context, id string, updates bson.M) (*mongo.UpdateResult, error) {
	update := bson.M{"$set": updates}
	result, err := repo.collection().UpdateOne(ctx, bson.M{"id": id}, update)
	return result, err
}

// DeleteEntity removes an entity from MongoDB
func (repo *MongoRepository) DeleteEntity(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	result, err := repo.collection().DeleteOne(ctx, bson.M{"id": id})
	return result, err
}
