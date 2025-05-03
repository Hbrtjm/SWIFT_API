package service

import (
	"fmt"
	"strings"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/models"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/repository"
	"github.com/Hbrtjm/SWIFT_API/backend/pkg/validators"
)

// PostBankData creates a new bank entry in the database
func (s *SwiftCodeService) PostBankData(bankData map[string]interface{}) error {
	bankValidator := validators.NewBankRequestValidator()
	err := bankValidator.ValidateAndSanitize(bankData)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	countryISO2, _ := bankData["countryISO2"].(string)
	countryName, _ := bankData["countryName"].(string)
	timeZone, _ := bankData["timeZone"].(string)

	country := models.Country{
		CountryISO2: strings.ToUpper(countryISO2),
		CountryName: strings.ToUpper(countryName),
		TimeZone:    timeZone,
	}

	// Add or update the country in the database
	err = s.repo.InsertCountry(country)
	if err != nil {
		return fmt.Errorf("failed to process country data: %w", err)
	}

	swiftCode, _ := bankData["swiftCode"].(string)

	bank := models.Bank{
		CountryISO2:   strings.ToUpper(countryISO2),
		SwiftCode:     swiftCode,
		CodeType:      getValue(bankData, "codeType"),
		BankName:      getValue(bankData, "bankName"),
		Address:       getValue(bankData, "address"),
		TownName:      getValue(bankData, "townName"),
		IsHeadquarter: getBool(bankData, "isHeadquarter"),
		BranchCode:    "",
	}

	// Get the first 8 characters of the SWFIT code for the branch code
	bank.BranchCode = swiftCode[:8]

	err = s.repo.InsertBank(bank)
	if err != nil {
		if err == repository.ErrBankExists {
			return fmt.Errorf("bank with SWIFT code %s already exists", swiftCode)
		}
		return fmt.Errorf("failed to insert bank: %w", err)
	}

	return nil
}
