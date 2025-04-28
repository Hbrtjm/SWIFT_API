package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FindBySwiftCode finds a document by SWIFT code
func (r *MongoRepository) FindBySwiftCode(code string) (map[string]interface{}, error) {
	if r == nil {
		return nil, errors.New("repository is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Filter by SWIFT code
	filter := bson.D{{Key: "swiftCode", Value: code}}

	// Get one result, if exists, otherwise return error
	var result models.Bank
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no bank found with SWIFT code %s", code)
		}
		return nil, err
	}

	return StructToMap(result)
}

func (r *MongoRepository) FindByBranchCode(branchCode string) ([]map[string]interface{}, error) {
	if r == nil {
		return nil, errors.New("repository is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Filter by branchCode field
	filter := bson.D{{Key: "branchCode", Value: branchCode}}

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
	fmt.Println("Banks found:", len(banks))
	fmt.Println("Banks:", banks)
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

// StructToMap converts a struct to a map using JSON marshaling
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
