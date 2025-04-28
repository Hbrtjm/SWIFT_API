package validators

import (
	"errors"
)

type SwiftCodeRequest struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

func ValidateRequest(request SwiftCodeRequest) error {
	// Validate required fields
	if request.SwiftCode == "" || request.BankName == "" || request.CountryISO2 == "" {
		return errors.New("missing required fields, swiftCode, bankName or countryISO2")
	}

	validator := SwiftCodeValidator{}
	if err := validator.Validate(request); err != nil {
		return err
	}

	if err := validator.ValidateWithCountryCode(request.SwiftCode, request.CountryISO2); err != nil {
		return err
	}

	return nil
}
