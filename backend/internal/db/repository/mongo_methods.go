package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/models"
	"go.mongodb.org/mongo-driver/bson"
)

// Delete deletes a document by SWIFT code
func (r *MongoRepository) Delete(code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	filter := bson.M{"swiftCode": code}
	result, err := r.bankCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("database delete error: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%s SWIFT code not found", code)
	}

	return nil
}

// InsertBank inserts a new bank into the database
func (r *MongoRepository) InsertBank(bank models.Bank) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	// Check if the bank already exists
	_, err := r.FindBySwiftCode(bank.SwiftCode)
	if err == nil {
		return ErrBankExists
	}

	// Insert the bank
	_, err = r.bankCollection.InsertOne(ctx, bank)
	if err != nil {
		return fmt.Errorf("database error during bank insertion: %w", err)
	}

	return nil
}

// InsertMany inserts multiple bank documents into the collection
func (r *MongoRepository) InsertManyBanks(banks []models.Bank) error {
	if len(banks) == 0 {
		return nil
	}

	// Convert to interface slice
	data := make([]interface{}, len(banks))
	for i := range banks {
		data[i] = banks[i]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.bankCollection.InsertMany(ctx, data)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

// InsertManyCountries inserts multiple country documents into the collection
func (r *MongoRepository) InsertManyCountries(countries []models.Country) error {
	if len(countries) == 0 {
		return nil
	}

	// Convert to interface slice
	data := make([]interface{}, len(countries))
	for i := range countries {
		data[i] = countries[i]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.countryCollection.InsertMany(ctx, data)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

// InsertCountry inserts a new country if it doesn't exist
func (r *MongoRepository) InsertCountry(country models.Country) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	exists, err := r.CountryExists(country.CountryISO2)
	if err != nil {
		return err
	}

	if exists {
		existingCountry, err := r.GetCountry(country.CountryISO2)
		if err != nil {
			return err
		}

		// If country name is different and the new one is not empty, update it
		if country.CountryName != "" && country.CountryName != existingCountry.CountryName {
			filter := bson.M{"countryISO2": country.CountryISO2}
			update := bson.M{"$set": bson.M{"countryName": country.CountryName}}

			_, err := r.countryCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				return fmt.Errorf("database error updating country: %w", err)
			}
		}

		return nil // Country exists, no need to insert
	}

	_, err = r.countryCollection.InsertOne(ctx, country)
	if err != nil {
		return fmt.Errorf("database error inserting country: %w", err)
	}

	return nil
}
