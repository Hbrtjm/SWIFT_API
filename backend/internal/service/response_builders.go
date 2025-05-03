package service

import (
	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/models"
)

// Helper function to map a Bank model to SwiftCodeResponse
func bankToResponse(bank *models.Bank, countryName string) *SwiftCodeResponse {
	return &SwiftCodeResponse{
		Address:       bank.Address,
		BankName:      bank.BankName,
		CountryISO2:   bank.CountryISO2,
		CountryName:   countryName,
		IsHeadquarter: bank.IsHeadquarter,
		SwiftCode:     bank.SwiftCode,
	}
}

// Helper function to map a Bank model to a map[string]interface{}
func mapBankToMap(bank *models.Bank) map[string]interface{} {
	return map[string]interface{}{
		"address":       bank.Address,
		"bankName":      bank.BankName,
		"countryISO2":   bank.CountryISO2,
		"isHeadquarter": bank.IsHeadquarter,
		"swiftCode":     bank.SwiftCode,
	}
}

func mapBranchValues(value map[string]interface{}) map[string]interface{} {
	address, _ := value["address"].(string)
	bankName, _ := value["bankName"].(string)
	countryISO2, _ := value["countryISO2"].(string)
	isHeadquarter, _ := value["isHeadquarter"].(bool)
	swiftCode, _ := value["swiftCode"].(string)

	return map[string]interface{}{
		"address":       address,
		"bankName":      bankName,
		"countryISO2":   countryISO2,
		"isHeadquarter": isHeadquarter,
		"swiftCode":     swiftCode,
	}
}
