package repository

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIndices creates database indices for better performance
func (r *MongoRepository) CreateIndices(logger *log.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create index on swiftCode field which is unique
	_, err := r.Collection().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "swiftCode", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		logger.Printf("Error creating swiftCode index: %v", err)
		return err
	}

	// Create index on countryCode field to optimize country code queries
	_, err = r.Collection().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "countryCode", Value: 1}},
	})
	if err != nil {
		logger.Printf("Error creating countryCode index: %v", err)
		return err
	}

	_, err = r.Collection().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "branchCode", Value: 1}},
	})

	// Create index on branchCode field to optimize branches searches
	if err != nil {
		logger.Printf("Error creating breanchCode index: %v", err)
		return err
	}

	logger.Println("Successfully created database indices")
	return nil
}