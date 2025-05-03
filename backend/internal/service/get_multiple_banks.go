package service

import (
	"errors"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/models"
)

// GetMultipleSwiftCodes returns data for multiple SWIFT codes
func (s *SwiftCodeService) GetMultipleSwiftCodes(codes []string) ([]map[string]interface{}, error) {
	if len(codes) == 0 {
		return []map[string]interface{}{}, nil
	}

	result := make([]map[string]interface{}, 0, len(codes))
	for _, code := range codes {
		bank, err := s.repo.FindBySwiftCode(code)
		emptyBank := models.Bank{}
		if err == nil && bank != emptyBank {
			bankMap := mapBankToMap(&bank)

			// Try to get country name, but don't fail if it's not found
			countryName, err := s.repo.LookupCountryName(bank.CountryISO2)
			if err == nil {
				bankMap["countryName"] = countryName
			}

			result = append(result, bankMap)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no valid SWIFT codes found")
	}

	return result, nil
}
