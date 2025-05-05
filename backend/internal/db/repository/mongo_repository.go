package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/api/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define custom errors for better error handling
var (
	ErrCountryExists   = errors.New("country already exists")
	ErrCountryNotFound = errors.New("country not found")
	ErrBankExists      = errors.New("bank already exists")
)

// MongoRepository handles database operations
type MongoRepository struct {
	client            *mongo.Client
	database          *mongo.Database
	bankCollection    *mongo.Collection
	countryCollection *mongo.Collection
}

// NewMongoRepository creates a new MongoRepository instance
func NewMongoRepository(uri, dbName, bankCollectionName, countriesCollectionName string, logger *middleware.Logger) (*MongoRepository, error) {
	var client *mongo.Client
	var err error

	maxRetries := 10
	delay := 5 * time.Second

	// Repeat connection retries until we connect, this will run for 50 seconds in the current setting
	for i := 0; i < maxRetries; i++ {
		client, err = NewMongoClient(uri)
		if err == nil && PingMongo(client) == nil {
			break
		}
		logger.Info("Mongo is not ready, retrying for the %d time out of %d", i, maxRetries)
		time.Sleep(delay)
	}

	// Connection failed
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w after %d retries", err, maxRetries)
	}

	// Check connection
	if err := PingMongo(client); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := GetMongoDatabase(client, dbName)
	bankCollection := GetMongoCollection(db, bankCollectionName)
	countriesCollection := GetMongoCollection(db, countriesCollectionName)

	return &MongoRepository{
		client:            client,
		database:          db,
		bankCollection:    bankCollection,
		countryCollection: countriesCollection,
	}, nil
}

func NewMongoClient(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri).SetMaxPoolSize(2000).SetMinPoolSize(10).SetConnectTimeout(20 * time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func PingMongo(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	return client.Ping(ctx, nil)
}

func GetMongoDatabase(client *mongo.Client, dbName string) *mongo.Database {
	return client.Database(dbName)
}

func GetMongoCollection(db *mongo.Database, collectionName string) *mongo.Collection {
	return db.Collection(collectionName)
}

func (r *MongoRepository) CountriesCollection() *mongo.Collection {
	return r.countryCollection
}

func (r *MongoRepository) BanksCollection() *mongo.Collection {
	return r.bankCollection
}

func (r *MongoRepository) CloseConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	if r.client != nil {
		return r.client.Disconnect(ctx)
	}
	return nil
}
