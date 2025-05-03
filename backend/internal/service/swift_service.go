package service

import (
	"errors"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/api/middleware"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/repository"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/parser"
)

// SwiftCodeService handles business logic for SWIFT codes
type SwiftCodeService struct {
	repo   *repository.MongoRepository
	parser *parser.SwiftFileParser
	logger *middleware.Logger
}

// NewSwiftCodeService creates a new SwiftCodeService
func NewSwiftCodeService(repo *repository.MongoRepository, parser *parser.SwiftFileParser, logger *middleware.Logger) *SwiftCodeService {
	return &SwiftCodeService{
		repo:   repo,
		parser: parser,
		logger: logger,
	}
}

// LoadInitialData parses and loads bank and country data from a file into the database
func (s *SwiftCodeService) LoadInitialData(filename string) error {
	// Parse the file into Bank and Country models
	banks, countries, err := s.parser.ParseFile(filename)
	if err != nil {
		s.logger.Error("Error parsing file: %v", err)
		return err
	}

	s.logger.Info("Parsed %d banks and %d countries from file", len(banks), len(countries))
	if len(banks) == 0 {
		return errors.New("no banks found in file")
	}

	// Insert the banks into MongoDB
	s.logger.Info("Inserting %d banks into database", len(banks))
	err = s.repo.InsertManyBanks(banks)
	if err != nil {
		return err
	}

	// Insert the countries into MongoDB
	s.logger.Info("Inserting %d countries into database", len(countries))
	err = s.repo.InsertManyCountries(countries)
	if err != nil {
		return err
	}

	return nil
}
