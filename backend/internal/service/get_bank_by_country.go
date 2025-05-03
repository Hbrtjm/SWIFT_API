package service

import (
	"fmt"
	"strings"

	"github.com/Hbrtjm/SWIFT_API/backend/pkg/validators"
)

// GetBySwiftCodesByCountry returns all SWIFT codes for a given country
func (s *SwiftCodeService) GetBySwiftCodesByCountry(countryISO2 string) (map[string]interface{}, error) {
	countryISO2Validator := validators.NewCountryISO2CodeValidator()
	if err := countryISO2Validator.Validate(countryISO2); err != nil {
		return nil, fmt.Errorf("invalid country code %s: %v", countryISO2, err)
	}

	countryISO2 = strings.ToUpper(countryISO2)
	countryName, err := s.repo.LookupCountryName(countryISO2)
	if err != nil {
		return nil, fmt.Errorf("country lookup failed: %w", err)
	}

	banks, err := s.repo.FindByCountry(countryISO2)
	if err != nil {
		return nil, fmt.Errorf("bank lookup failed: %w", err)
	}

	// Map bank data to response format
	swiftCodes := make([]map[string]interface{}, 0, len(banks))
	for _, bank := range banks {
		swiftCodes = append(swiftCodes, mapBankToMap(&bank))
	}

	response := map[string]interface{}{
		"countryISO2": countryISO2,
		"countryName": countryName,
		"swiftCodes":  swiftCodes,
	}

	return response, nil
}
