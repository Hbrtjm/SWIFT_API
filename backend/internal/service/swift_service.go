package service

import (
	"errors"
	"log"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/repository"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/parser"
)

// SwiftCodeService handles business logic for SWIFT codes
type SwiftCodeService struct {
	repo   *repository.MongoRepository
	parser *parser.SwiftFileParser
	logger *log.Logger
}

// NewSwiftCodeService creates a new SwiftCodeService
func NewSwiftCodeService(repo *repository.MongoRepository, parser *parser.SwiftFileParser, logger *log.Logger) *SwiftCodeService {
	return &SwiftCodeService{
		repo:   repo,
		parser: parser,
		logger: logger,
	}
}

// LoadInitialData parses and loads bank data from a file into the database
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
		// Extract branch code (first 8 characters)
		if len(bank.SwiftCode) >= 8 {
			bank.BranchCode = bank.SwiftCode[:8]
		}
		data[i] = bank
	}

	// Insert the data into MongoDB
	s.logger.Printf("Inserting %d banks into database", len(data))
	// Swich to separate inserts
	return s.repo.InsertMany(data)
}
