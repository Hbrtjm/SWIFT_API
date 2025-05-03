package repository

import (
	"context"
	"time"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/api/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIndices creates database indices for better performance
func (r *MongoRepository) CreateIndices(logger *middleware.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create index on swiftCode field which is unique
	_, err := r.BanksCollection().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "swiftCode", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		logger.Error("Error creating swiftCode index in banks collection: %v", err)
		return err
	}

	// Create index on countryISO2 field to optimize country code queries
	_, err = r.BanksCollection().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "countryISO2", Value: 1}},
	})
	if err != nil {
		logger.Error("Error creating countryISO2 index in banks collection: %v", err)
		return err
	}

	_, err = r.BanksCollection().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "branchCode", Value: 1}},
	})

	// Create index on branchCode field to optimize branches searches
	if err != nil {
		logger.Error("Error creating branchCode index in banks collection: %v", err)
		return err
	}

	_, err = r.CountriesCollection().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "countryISO2", Value: 1}},
	})

	if err != nil {
		logger.Error("Error creating countryISO2 index in countries collection: %v", err)
		return err
	}

	logger.Info("Successfully created database indices")
	return nil
}
