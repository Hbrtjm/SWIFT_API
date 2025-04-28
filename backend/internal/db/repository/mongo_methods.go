package repository

import (
	"context"
	"time"
	"errors"
	"fmt"
	"strings"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/models"
	"go.mongodb.org/mongo-driver/bson"
)

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