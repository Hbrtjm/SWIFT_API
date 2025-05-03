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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindBySwiftCode finds a document by SWIFT code
// FindBySwiftCode finds a bank by SWIFT code
func (r *MongoRepository) FindBySwiftCode(swiftCode string) (models.Bank, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	var bank models.Bank
	filter := bson.M{"swiftCode": swiftCode}
	err := r.bankCollection.FindOne(ctx, filter).Decode(&bank)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Bank{}, fmt.Errorf("no bank found with SWIFT code %s", swiftCode)
		}
		return models.Bank{}, err
	}
	return bank, nil
}

func (r *MongoRepository) FindByBranchCode(branchCode string) ([]map[string]interface{}, error) {
	if r == nil {
		return nil, errors.New("repository is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Filter by branchCode field
	filter := bson.D{{Key: "branchCode", Value: branchCode}}

	cursor, err := r.bankCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// First decode into a slice of Bank structs
	var banks []models.Bank
	if err = cursor.All(ctx, &banks); err != nil {
		return nil, err
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

// FindByCountry finds all banks by country ISO2 code
func (r *MongoRepository) FindByCountry(countryISO2 string) ([]models.Bank, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	var banks []models.Bank
	filter := bson.M{"countryISO2": countryISO2}

	// Todo - can be deleted
	opts := options.Find().SetSort(bson.D{{Key: "swiftCode", Value: 1}})

	cursor, err := r.bankCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &banks); err != nil {
		return nil, err
	}

	if len(banks) == 0 {
		return nil, fmt.Errorf("no banks found for country %s", countryISO2)
	}

	return banks, nil
}

// Count returns the total number of documents in the collection
func (r *MongoRepository) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	count, err := r.bankCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetCountry retrieves a country by ISO2 code
func (r *MongoRepository) GetCountry(countryISO2 string) (models.Country, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	var country models.Country
	filter := bson.M{"countryISO2": countryISO2}

	err := r.countryCollection.FindOne(ctx, filter).Decode(&country)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Country{}, ErrCountryNotFound
		}
		return models.Country{}, fmt.Errorf("database error retrieving country: %w", err)
	}

	return country, nil
}

// LookupCountryName looks up a country name by ISO2 code
func (r *MongoRepository) LookupCountryName(countryISO2 string) (string, error) {
	country, err := r.GetCountry(countryISO2)
	if err != nil {
		return "", err
	}
	return country.CountryName, nil
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

// CountryExists checks if a country exists in the database
func (r *MongoRepository) CountryExists(countryISO2 string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	filter := bson.M{"countryISO2": countryISO2}
	count, err := r.countryCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("database error checking country existence: %w", err)
	}

	return count > 0, nil
}
