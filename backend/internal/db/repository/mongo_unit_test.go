package repository

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Hbrtjm/SWIFT_API/internal/db/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	testMongoURI = "mongodb://localhost:27017"
	testDBName   = "swiftcodes_test"
)

var repo *MongoRepository

// Setup function to initialize the test database
func setupTestDB() (*MongoRepository, error) {
	r, err := NewMongoRepository(testMongoURI, testDBName)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = r.database.Drop(ctx)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// TestMain runs before all tests
func TestMain(m *testing.M) {
	var err error
	repo, err = setupTestDB()
	if err != nil {
		log.Fatalf("Failed to set up test database: %v", err)
	}
	code := m.Run()

	if repo != nil {
		repo.Close()
	}

	os.Exit(code)
}

// addTestData adds test banks to the database
func addTestData(t *testing.T) {
	testBanks := []models.Bank{
		{
			CountryCode:  "AL",
			SwiftCode:    "AAISALTRXXX",
			CodeType:     "BIC11",
			BankName:     "UNITED BANK OF ALBANIA SH.A",
			Address:      "HYRJA 3 RR. DRITAN HOXHA ND. 11 TIRANA, TIRANA, 1023",
			TownName:     "TIRANA",
			CountryName:  "ALBANIA",
			TimeZone:     "Europe/Tirane",
			IsHeadOffice: true,
			BranchCode:   "XXX",
		},
		{
			CountryCode:  "BG",
			SwiftCode:    "ABIEBGS1XXX",
			CodeType:     "BIC11",
			BankName:     "ABV INVESTMENTS LTD",
			Address:      "TSAR ASEN 20  VARNA, VARNA, 9002",
			TownName:     "VARNA",
			CountryName:  "BULGARIA",
			TimeZone:     "Europe/Sofia",
			IsHeadOffice: true,
			BranchCode:   "XXX",
		},
		{
			CountryCode:  "UY",
			SwiftCode:    "AFAAUYM1XXX",
			CodeType:     "BIC11",
			BankName:     "AFINIDAD A.F.A.P.S.A.",
			Address:      "PLAZA INDEPENDENCIA 743  MONTEVIDEO, MONTEVIDEO, 11000",
			TownName:     "MONTEVIDEO",
			CountryName:  "URUGUAY",
			TimeZone:     "America/Montevideo",
			IsHeadOffice: true,
			BranchCode:   "XXX",
		},
	}

	for _, bank := range testBanks {
		err := repo.Insert(bank)
		require.NoError(t, err)
	}
}

// cleanTestData removes all test data
func cleanTestData(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.collection.Drop(ctx)
	require.NoError(t, err)
}

// TestNewMongoRepository tests the creation of a new repository
func TestNewMongoRepository(t *testing.T) {
	repo, err := NewMongoRepository(testMongoURI, testDBName)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	repo, err = NewMongoRepository("mongodb://invalid:27017", testDBName)
	assert.Error(t, err)
	assert.Nil(t, repo)
}

// TestInsert tests inserting a single document
func TestInsert(t *testing.T) {
	cleanTestData(t)

	bank := models.Bank{
		CountryCode:  "AL",
		SwiftCode:    "AAISALTRXXX",
		CodeType:     "BIC11",
		BankName:     "UNITED BANK OF ALBANIA SH.A",
		Address:      "HYRJA 3 RR. DRITAN HOXHA ND. 11 TIRANA, TIRANA, 1023",
		TownName:     "TIRANA",
		CountryName:  "ALBANIA",
		TimeZone:     "Europe/Tirane",
		IsHeadOffice: true,
		BranchCode:   "XXX",
	}
	err := repo.Insert(bank)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result models.Bank
	err = repo.collection.FindOne(ctx, bson.M{"swiftCode": "AAISALTRXXX"}).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, bank.BankName, result.BankName)
	assert.Equal(t, bank.Address, result.Address)
	assert.Equal(t, bank.TimeZone, result.TimeZone)
}

// TestInsertMany tests inserting multiple documents
func TestInsertMany(t *testing.T) {
	cleanTestData(t)

	banks := []interface{}{
		models.Bank{
			CountryCode:  "AL",
			SwiftCode:    "MULTAL123XXX",
			CodeType:     "BIC11",
			BankName:     "Multi Bank Albania",
			Address:      "Street 1, Tirana",
			TownName:     "TIRANA",
			CountryName:  "ALBANIA",
			TimeZone:     "Europe/Tirane",
			IsHeadOffice: true,
			BranchCode:   "XXX",
		},
		models.Bank{
			CountryCode:  "UY",
			SwiftCode:    "MULTUY456XXX",
			CodeType:     "BIC11",
			BankName:     "Multi Bank Uruguay",
			Address:      "Street 2, Montevideo",
			TownName:     "MONTEVIDEO",
			CountryName:  "URUGUAY",
			TimeZone:     "America/Montevideo",
			IsHeadOffice: true,
			BranchCode:   "XXX",
		},
	}

	err := repo.InsertMany(banks)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := repo.collection.CountDocuments(ctx, bson.M{})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

// TestFindBySwiftCode tests finding a document by SWIFT code
func TestFindBySwiftCode(t *testing.T) {
	cleanTestData(t)
	addTestData(t)

	result, err := repo.FindBySwiftCode("AAISALTRXXX")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "AAISALTRXXX", result["swiftCode"])
	assert.Equal(t, "UNITED BANK OF ALBANIA SH.A", result["bankName"])
	assert.Equal(t, "TIRANA", result["townName"])
	assert.Equal(t, "Europe/Tirane", result["timeZone"])
	assert.Equal(t, "BIC11", result["codeType"])

	result, err = repo.FindBySwiftCode("NONEXISTENTOOO")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no bank found with SWIFT code NONEXISTENTOOO")
}

// TestFindByCountry tests finding documents by country code
func TestFindByCountry(t *testing.T) {
	cleanTestData(t)
	addTestData(t)

	results, err := repo.FindByCountry("BG")
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "ABIEBGS1XXX", results[0]["swiftCode"])

	results, err = repo.FindByCountry("XX")
	assert.Error(t, err)
	assert.Nil(t, results)
}

// TestCount tests counting documents
func TestCount(t *testing.T) {
	cleanTestData(t)

	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	addTestData(t)

	count, err = repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

// TestDelete tests deleting a document
func TestDelete(t *testing.T) {
	cleanTestData(t)
	addTestData(t)

	err := repo.Delete("AAISALTRXXX")
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := repo.collection.CountDocuments(ctx, bson.M{})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)

	// Test deleting non-existing document
	err = repo.Delete("NONEXISTENTOOO")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
