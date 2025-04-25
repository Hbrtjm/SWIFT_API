package service

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/Hbrtjm/SWIFT_API/internal/db/models"
	"github.com/Hbrtjm/SWIFT_API/internal/db/repository"
	"github.com/Hbrtjm/SWIFT_API/internal/parser"
)

// SwiftCodeService handles business logic for SWIFT codes
type SwiftCodeService struct {
	repo   *repository.MongoRepository
	parser *parser.SwiftCodeParser
	logger *log.Logger
}

// NewSwiftCodeService creates a new SwiftCodeService
func NewSwiftCodeService(repo *repository.MongoRepository, parser *parser.SwiftCodeParser, logger *log.Logger) *SwiftCodeService {
	return &SwiftCodeService{
		repo:   repo,
		parser: parser,
		logger: logger,
	}
}

func (s *SwiftCodeService) LoadInitialData(filename string) error {
	// Parse the file into Bank models
	banks, err := s.parser.ParseFile(filename)
	if err != nil {
		s.logger.Printf("Error parsing file: %v", err)
		return err
	}

	s.logger.Printf("Parsed %d banks from file", len(banks))
	if len(banks) == 0 {
		return errors.New("no banks found in file")
	}

	// Convert banks from the file to the interface array
	data := make([]interface{}, len(banks))
	for i, bank := range banks {
		// TODO - SWIFT code validation before inserting
		data[i] = bank
	}

	// Insert the data into MongoDB
	s.logger.Printf("Inserting %d banks into database", len(data))
	return s.repo.InsertMany(data)
}

// GetSwiftCode returns a SWIFT code by its identifier
func (s *SwiftCodeService) GetSwiftCode(code string) (map[string]interface{}, error) {

	// Query the repository
	return s.repo.FindBySwiftCode(code)
}

// GetSwiftCodesByCountry returns all SWIFT codes for a given country
func (s *SwiftCodeService) GetSwiftCodesByCountry(countryCode string) ([]map[string]interface{}, error) {
	if len(countryCode) != 2 {
		return nil, errors.New("country code must be ISO-2 format (2 characters)")
	}

	// Convert to uppercase to match the database format
	countryCode = strings.ToUpper(countryCode)

	return s.repo.FindByCountry(countryCode)
}

// DeleteSwiftCode deletes a SWIFT code
func (s *SwiftCodeService) DeleteSwiftCode(code string) error {
	// Delete from the repository
	return s.repo.Delete(code)
}

// GetMultipleSwiftCodes returns data for multiple SWIFT codes
func (s *SwiftCodeService) GetMultipleSwiftCodes(codes []string) ([]map[string]interface{}, error) {
	if len(codes) == 0 {
		return []map[string]interface{}{}, nil
	}

	result := make([]map[string]interface{}, 0, len(codes))
	for _, code := range codes {
		data, err := s.repo.FindBySwiftCode(code)
		if err == nil {
			result = append(result, data)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no valid SWIFT codes found")
	}

	return result, nil
}

func (s *SwiftCodeService) CreateSwiftCode(bankData map[string]interface{}) error {
	bankJson, err := json.Marshal(bankData)
	if err != nil {
		return err
	}

	var bank models.Bank
	err = json.Unmarshal(bankJson, &bank)
	if err != nil {
		return err
	}

	return s.repo.Insert(bank)
}
