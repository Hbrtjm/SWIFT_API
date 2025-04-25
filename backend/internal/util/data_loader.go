package util

import (
	"context"
	"log"
	"time"

	"github.com/Hbrtjm/SWIFT_API/internal/db/repository"
	"github.com/Hbrtjm/SWIFT_API/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LoadInitialDataIfNeeded checks if database is empty and loads initial data if needed
func LoadInitialDataIfNeeded(service *service.SwiftCodeService, repo *repository.MongoRepository, filename string, logger *log.Logger) error {
	// TODO - Check if all of the data is already in the database, compare the sets
	count, err := repo.Count()
	if err != nil {
		logger.Printf("Error checking database: %v", err)
		return err
	}

	// If we already have data, skip initial data load
	if count > 0 {
		logger.Printf("Database already contains %d SWIFT codes, skipping initial data load", count)
		return nil
	}

	logger.Printf("Database is empty. Loading initial SWIFT data from %s", filename)
	err = service.LoadInitialData(filename)
	if err != nil {
		logger.Printf("Error loading initial data: %v", err)
		return err
	}

	newRowCount, _ := repo.Count()
	logger.Printf("Successfully loaded %d SWIFT codes into the database", newRowCount)

	return nil
}

// CreateIndices creates database indices for better performance
func CreateIndices(repo *repository.MongoRepository, logger *log.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create index on swiftCode field which is unique
	_, err := repo.Collection().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "swiftCode", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		logger.Printf("Error creating swiftCode index: %v", err)
		return err
	}

	// Create index on countryCode field to optimize country code queries
	_, err = repo.Collection().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "countryCode", Value: 1}},
	})
	if err != nil {
		logger.Printf("Error creating countryCode index: %v", err)
		return err
	}

	logger.Println("Successfully created database indices")
	return nil
}
