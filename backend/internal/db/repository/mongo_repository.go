package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository handles MongoDB operations
type MongoRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

// NewMongoRepository creates a new MongoRepository instance
func NewMongoRepository(uri, dbName, collectionName string) (*MongoRepository, error) {
	client, err := NewMongoClient(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	if err := PingMongo(client); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := GetMongoDatabase(client, dbName)
	coll := GetMongoCollection(db, collectionName)

	return &MongoRepository{
		client:     client,
		database:   db,
		collection: coll,
	}, nil
}

func NewMongoClient(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func PingMongo(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return client.Ping(ctx, nil)
}

func GetMongoDatabase(client *mongo.Client, dbName string) *mongo.Database {
	return client.Database(dbName)
}

func GetMongoCollection(db *mongo.Database, collectionName string) *mongo.Collection {
	return db.Collection(collectionName)
}

func (r *MongoRepository) Collection() *mongo.Collection {
	return r.collection
}

// Close closes the MongoDB connection
func (r *MongoRepository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if r.client != nil {
		return r.client.Disconnect(ctx)
	}
	return nil
}
