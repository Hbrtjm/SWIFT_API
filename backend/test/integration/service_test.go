package integration

import (
	"context"
	"log"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/api/middleware"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/models"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/repository"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/parser"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testMongoURI                = "mongodb://localhost:27017"
	testDBName                  = "swiftcodes_integration_test"
	testBankCollectionName      = "banks_integration_test"
	testCountriesCollectionName = "countries_integration_test"
	testDataFilePath            = "testdata/sample_swift_codes.csv"
)

var (
	repo         *repository.MongoRepository
	swiftParser  *parser.SwiftFileParser
	swiftService *service.SwiftCodeService
	testLogger   *middleware.Logger
)

// setupTestEnvironment initializes the test environment
func setupTestEnvironment() error {
	var err error

	// Create a logger that writes to a temporary file
	logFile, err := os.CreateTemp("", "swift_service_test_*.log")
	if err != nil {
		return err
	}

	testLogger = middleware.New(logFile, "", true)

	repo, err = repository.NewMongoRepository(testMongoURI, testDBName, testBankCollectionName, testCountriesCollectionName, middleware.NewNoLogger())
	if err != nil {
		return err
	}
	swiftParser = parser.NewSwiftFileParser()
	swiftService = service.NewSwiftCodeService(repo, swiftParser, testLogger)

	return nil
}

// cleanupTestDB drops the test collection
func cleanupTestDB() {
	if repo != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancel()

		repo.BanksCollection().Drop(ctx)
	}
}

// TestMain runs before all tests
func TestMain(m *testing.M) {
	// Set up the test environment
	if err := setupTestEnvironment(); err != nil {
		log.Fatalf("Failed to set up test environment: %v", err)
	}

	createSampleDataFile()

	code := m.Run()

	cleanupTestDB()
	if repo != nil {
		repo.CloseConnection()
	}

	os.Remove(testDataFilePath)

	os.Exit(code)
}

// createSampleDataFile creates a CSV file with sample data for testing
func createSampleDataFile() {
	sampleData := `COUNTRY ISO2 CODE;SWIFT CODE;CODE TYPE;NAME;ADDRESS;TOWN NAME;COUNTRY NAME;TIME ZONE
PL;TPEOPLPWXXX;BIC11;PEKAO TFI S.A.;FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066;WARSZAWA;POLAND;Europe/Warsaw
PL;TPEOPLPWP65;BIC11;PEKAO TOWARZYSTWO FUNDUSZY  INWESTYCYJNYCH SPOLKA AKCYJNA;FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066;WARSZAWA;POLAND;Europe/Warsaw
PL;TPEOPLPWPAE;BIC11;PEKAO TOWARZYSTWO FUNDUSZY  INWESTYCYJNYCH SPOLKA AKCYJNA;FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066;WARSZAWA;POLAND;Europe/Warsaw
PL;TPEOPLPWPFI;BIC11;PEKAO TOWARZYSTWO FUNDUSZY  INWESTYCYJNYCH SPOLKA AKCYJNA;FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066;WARSZAWA;POLAND;Europe/Warsaw
`
	err := os.MkdirAll("testdata", 0755)
	if err != nil {
		log.Fatalf("Failed to create testdata directory: %v", err)
	}

	err = os.WriteFile(testDataFilePath, []byte(sampleData), 0644)
	if err != nil {
		log.Fatalf("Failed to create sample data file: %v", err)
	}
}

// cleanup drops the collection between tests
func cleanup(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := repo.BanksCollection().Drop(ctx)
	require.NoError(t, err, "Failed to drop collection")
}

// TestLoadInitialData tests the LoadInitialData function
func TestLoadInitialData(t *testing.T) {
	cleanup(t)

	// Test loading data from file
	err := swiftService.LoadInitialData(testDataFilePath)
	assert.NoError(t, err)

	// Verify that data was correctly loaded
	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(4), count, "Expected 4 banks to be loaded")

	bank, err := repo.FindBySwiftCode("TPEOPLPWXXX")
	assert.NoError(t, err)
	assert.Equal(t, "PEKAO TFI S.A.", bank.BankName)
	assert.Equal(t, "PL", bank.CountryISO2)

	// Test loading from non-existent file
	err = swiftService.LoadInitialData("non_existent_file.csv")
	assert.Error(t, err)
}

// TestPostBankData tests the PostBankData function
func TestPostBankData(t *testing.T) {
	cleanup(t)

	// Test posting a new bank
	bankData := map[string]interface{}{
		"countryISO2":   "LV",
		"swiftCode":     "AIZKLV22XXX",
		"codeType":      "BIC11",
		"bankName":      "ABLV BANK, AS IN LIQUIDATION",
		"address":       "MIHAILA TALA STREET 1  RIGA, RIGA, LV-1045",
		"townName":      "RIGA",
		"countryName":   "LATVIA",
		"isHeadquarter": true,
		"timeZone":      "Europe/Riga",
	}

	err := swiftService.PostBankData(bankData)
	assert.NoError(t, err)

	// Check if the bank was added
	bank, err := repo.FindBySwiftCode("AIZKLV22XXX")
	assert.NoError(t, err)
	assert.Equal(t, "ABLV BANK, AS IN LIQUIDATION", bank.BankName)
	value, _ := repo.LookupCountryName(bank.CountryISO2)
	assert.Equal(t, "LATVIA", value)

	err = swiftService.PostBankData(bankData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bank with SWIFT code AIZKLV22XXX already exists")

	invalidBankData := map[string]interface{}{
		"notAValidField": "test",
	}
	err = swiftService.PostBankData(invalidBankData)
	assert.Error(t, err)
}

// TestGetBySwiftCode tests the GetBySwiftCode function
func TestGetBySwiftCode(t *testing.T) {
	cleanup(t)

	err := swiftService.LoadInitialData(testDataFilePath)
	assert.NoError(t, err)

	response, err := swiftService.GetBySwiftCode("TPEOPLPWXXX")
	assert.NotNil(t, response)
	assert.NoError(t, err)
	assert.Equal(t, "PEKAO TFI S.A.", response.BankName)
	assert.Equal(t, "PL", response.CountryISO2)

	response, err = swiftService.GetBySwiftCode("NONEXISTENT")
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "no bank found with the given SWIFT code")

	// Create an array of expected branch data - sorted by swiftCode lexicographically (refer to the strings.Compare() definition)
	expectedBranches := []map[string]string{
		{
			"swiftCode":  "TPEOPLPWPFI",
			"bankName":   "PEKAO TOWARZYSTWO FUNDUSZY  INWESTYCYJNYCH SPOLKA AKCYJNA",
			"branchCode": "TPEOPLPW",
			"address":    "FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066",
			"townName":   "WARSZAWA",
		},
		{
			"swiftCode":  "TPEOPLPWPAE",
			"bankName":   "PEKAO TOWARZYSTWO FUNDUSZY  INWESTYCYJNYCH SPOLKA AKCYJNA",
			"branchCode": "TPEOPLPW",
			"address":    "FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066",
			"townName":   "WARSZAWA",
		},
		{
			"swiftCode":  "TPEOPLPWP65",
			"bankName":   "PEKAO TOWARZYSTWO FUNDUSZY  INWESTYCYJNYCH SPOLKA AKCYJNA",
			"branchCode": "TPEOPLPW",
			"address":    "FOREST ZUBRA 1, FLOOR 1 WARSZAWA, MAZOWIECKIE, 01-066",
			"townName":   "WARSZAWA",
		},
	}

	response, err = swiftService.GetBySwiftCode("TPEOPLPWXXX")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "PEKAO TFI S.A.", response.BankName)
	assert.True(t, response.IsHeadquarter)

	// The headquarter should have branches in its branches list
	assert.NotNil(t, response.Branches)
	assert.Len(t, response.Branches, len(expectedBranches))

	// Sort the response branches by SWIFT code to ensure consistent comparison
	sort.Slice(response.Branches, func(i, j int) bool {
		return strings.Compare(response.Branches[i]["swiftCode"].(string), response.Branches[j]["swiftCode"].(string)) > 0
	})

	for i, branch := range response.Branches {
		assert.Equal(t, expectedBranches[i]["swiftCode"], branch["swiftCode"], "Branch %d swift code mismatch", i)
		assert.Equal(t, expectedBranches[i]["bankName"], branch["bankName"], "Branch %d bank name mismatch", i)
		assert.Equal(t, expectedBranches[i]["address"], branch["address"], "Branch %d address mismatch", i)
	}
}

// TestGetMultipleSwiftCodes tests the GetMultipleSwiftCodes function
func TestGetMultipleSwiftCodes(t *testing.T) {
	cleanup(t)

	err := swiftService.LoadInitialData(testDataFilePath)
	assert.NoError(t, err)

	// Test getting multiple banks by SWIFT codes
	codes := []string{"TPEOPLPWXXX", "TPEOPLPWPAE", "NONEXISTENT"}
	response, err := swiftService.GetMultipleSwiftCodes(codes)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response, 2)

	var foundBankA, foundBankB bool
	for _, bank := range response {
		if bank["swiftCode"] == "TPEOPLPWXXX" {
			foundBankA = true
			assert.Equal(t, "PEKAO TFI S.A.", bank["bankName"])
		}
		if bank["swiftCode"] == "TPEOPLPWPAE" {
			foundBankB = true
			assert.Equal(t, "PEKAO TOWARZYSTWO FUNDUSZY  INWESTYCYJNYCH SPOLKA AKCYJNA", bank["bankName"])
		}
	}
	assert.True(t, foundBankA, "Expected to find TPEOPLPWXXX")
	assert.True(t, foundBankB, "Expected to find TPEOPLPWPAE")

	// Test with empty list
	response, err = swiftService.GetMultipleSwiftCodes([]string{})
	assert.NoError(t, err)
	assert.Empty(t, response)

	// A couple of non-existent codes
	response, err = swiftService.GetMultipleSwiftCodes([]string{"NONEXISTENT1", "NONEXISTENT2"})
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "no valid SWIFT codes found")
}

// TestGetBySwiftCodesByCountry tests the GetBySwiftCodesByCountry function
func TestGetBySwiftCodesByCountry(t *testing.T) {
	cleanup(t)

	// Load initial data
	err := swiftService.LoadInitialData(testDataFilePath)
	assert.NoError(t, err)

	// Test getting banks by country code
	response, err := swiftService.GetBySwiftCodesByCountry("PL")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response, 3)
	assert.Len(t, response["swiftCodes"], 4)

	response, err = swiftService.GetBySwiftCodesByCountry("XX")
	assert.Error(t, err)
	assert.Nil(t, response)

	response, err = swiftService.GetBySwiftCodesByCountry("INVALID")
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid country code")
}

// TestDeleteSwiftCode tests the DeleteSwiftCode function
func TestDeleteSwiftCode(t *testing.T) {
	cleanup(t)

	// Load initial data
	err := swiftService.LoadInitialData(testDataFilePath)
	assert.NoError(t, err)

	// Test deleting a bank by SWIFT code
	err = swiftService.DeleteSwiftCode("TPEOPLPWXXX")
	assert.NoError(t, err)

	bank, err := repo.FindBySwiftCode("TPEOPLPWXXX")
	assert.Error(t, err)
	emptyBank := models.Bank{}
	assert.Equal(t, emptyBank, bank)

	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)

	err = swiftService.DeleteSwiftCode("NONEXISTENT")
	assert.Error(t, err)

	err = swiftService.DeleteSwiftCode("INVALID")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid SWIFT code")
}

// TestIntegrationFlow tests the entire flow from loading data to querying and modifying
func TestIntegrationFlow(t *testing.T) {
	cleanup(t)

	// Load initial data
	err := swiftService.LoadInitialData(testDataFilePath)
	assert.NoError(t, err)

	// Add a new bank
	newBank := map[string]interface{}{
		"countryISO2":   "PL",
		"swiftCode":     "TPEOPLPW123",
		"codeType":      "BIC11",
		"bankName":      "PEKAO TOWARZYSTWO FUNDUSZY  INWESTYCYJNYCH SPOLKA AKCYJNA",
		"address":       "456 TEST STRASSE",
		"townName":      "WARSAW",
		"countryName":   "POLAND",
		"isHeadquarter": false,
	}
	err = swiftService.PostBankData(newBank)
	assert.NoError(t, err)

	// Query the new bank and verify
	response, err := swiftService.GetBySwiftCode("TPEOPLPW123")
	assert.NoError(t, err)
	assert.Equal(t, "PEKAO TOWARZYSTWO FUNDUSZY  INWESTYCYJNYCH SPOLKA AKCYJNA", response.BankName)
	assert.Equal(t, "PL", response.CountryISO2)

	// Query multiple codes
	multiResponse, err := swiftService.GetMultipleSwiftCodes([]string{"TPEOPLPW123", "TPEOPLPWXXX"})
	assert.NoError(t, err)
	assert.Len(t, multiResponse, 2)

	// Delete added bank
	err = swiftService.DeleteSwiftCode("TPEOPLPW123")
	assert.NoError(t, err)

	// Verify deletion
	response, err = swiftService.GetBySwiftCode("TPEOPLPW123")
	assert.Error(t, err)
	assert.Nil(t, response)

	// Query by country
	countryResponse, err := swiftService.GetBySwiftCodesByCountry("PL")
	assert.NoError(t, err)
	assert.Len(t, countryResponse, 3)

	// Final count verification
	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(4), count)
}
