// backend/internal/db/repository/mongo_repository.go

package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Hbrtjm/SWIFT_API/internal/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO - Match the hedquarter and branch codes

// MongoRepository handles MongoDB operations
type MongoRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

// Constructor for MongoRepository structure instance
func NewMongoRepository(uri, dbName string) (*MongoRepository, error) {
	// Set connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Get database and collection
	db := client.Database(dbName)
	collection := db.Collection("swiftCodes")

	return &MongoRepository{
		client:     client,
		database:   db,
		collection: collection,
	}, nil
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

// Insert inserts a new document into the collection
func (r *MongoRepository) Insert(bank models.Bank) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, errExists := r.FindBySwiftCode(bank.SwiftCode)
	if errExists == nil {
		return errors.New("SWIFT code already exists")
	} else {
		if !strings.Contains(errExists.Error(), "no bank found") {
			return errExists
		}
	}

	_, err := r.collection.InsertOne(ctx, bank)
	if err != nil {
		return err
	}
	return nil
}

// InsertMany inserts multiple documents into the collection
func (r *MongoRepository) InsertMany(data []interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.InsertMany(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

// FindBySwiftCode finds a document by SWIFT code
func (r *MongoRepository) FindBySwiftCode(code string) (map[string]interface{}, error) {
	if r == nil { // Previous error, I think I can remove this check
		return nil, errors.New("repository is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	fmt.Printf("DB: %s, Collection: %s\n", r.database.Name(), r.collection.Name())
	fmt.Printf("Looking for SWIFT code: %s\n", code)

	// Filter by SWIFT code
	filter := bson.D{{Key: "swiftCode", Value: code}}

	// Get one result, if exists, otherwise return error
	var result models.Bank
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Printf("Error finding SWIFT code: %v\n", err)
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no bank found with SWIFT code %s", code)
		}
		return nil, err
	}

	return StructToMap(result)
}

// FindByCountry finds documents by country code
func (r *MongoRepository) FindByCountry(countryCode string) ([]map[string]interface{}, error) {
	if r == nil {
		return nil, errors.New("repository is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Use proper BSON filter
	filter := bson.D{{Key: "countryCode", Value: countryCode}}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// First decode into a slice of Bank structs
	var banks []models.Bank
	if err = cursor.All(ctx, &banks); err != nil {
		return nil, err
	}

	if len(banks) == 0 {
		return nil, fmt.Errorf("no banks found for country code %s", countryCode)
	}

	// Then convert each Bank struct to a map
	results := make([]map[string]interface{}, len(banks))
	for i, bank := range banks {
		bankMap, err := StructToMap(bank)
		if err != nil {
			return nil, err
		}
		results[i] = bankMap
	}

	return results, nil
}

// Count returns the total number of documents in the collection
func (r *MongoRepository) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Delete deletes a document by SWIFT code
func (r *MongoRepository) Delete(code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"swiftCode": code}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%s SWIFT code not found", code)
	}

	return nil
}

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}

	// Marshal struct to JSON
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON into map
	err = json.Unmarshal(jsonBytes, &result)
	return result, err
}
