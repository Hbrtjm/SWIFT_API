package service

import (
	"fmt"
	"strings"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/models"
	"github.com/Hbrtjm/SWIFT_API/backend/pkg/validators"
)

// PostCountry creates a new country entry in the database
func (s *SwiftCodeService) PostCountry(countryData map[string]interface{}) error {
	countryISO2, ok := countryData["countryISO2"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid countryISO2")
	}

	validator := validators.NewCountryValidator()
	if err := validator.ValidateAndSanitize(countryData); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	countryName, _ := countryData["countryName"].(string)
	timeZone, _ := countryData["timeZone"].(string)

	country := models.Country{
		CountryISO2: strings.ToUpper(countryISO2),
		CountryName: strings.ToUpper(countryName),
		TimeZone:    timeZone,
	}

	err := s.repo.InsertCountry(country)
	if err != nil {
		return fmt.Errorf("failed to process country data: %w", err)
	}

	return nil
}
