package service

import (
	"fmt"
	"strings"

	"github.com/Hbrtjm/SWIFT_API/backend/pkg/validators"
)

// GetBySwiftCodesByCountry returns all SWIFT codes for a given country
func (s *SwiftCodeService) GetBySwiftCodesByCountry(countryCode string) ([]map[string]interface{}, error) {
	countryCodeValidator := validators.NewCountryCodeValidator()
	if err := countryCodeValidator.Validate(countryCode); err != nil {
		return nil, fmt.Errorf("invalid country code %s: %v", countryCode, err)
	}

	// Convert to uppercase to match the database format
	countryCode = strings.ToUpper(countryCode)

	return s.repo.FindByCountry(countryCode)
}
