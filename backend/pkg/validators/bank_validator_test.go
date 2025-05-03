package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBankRequestValidator_ValidInput(t *testing.T) {
	validator := NewBankRequestValidator()

	data := map[string]interface{}{
		"address":       "123 Main St",
		"bankName":      "Test Bank",
		"countryISO2":   "DE",
		"swiftCode":     "DEUTDEFF",
		"codeType":      "BIC8",
		"isHeadquarter": false,
		"timeZone":      "Europe/Berlin",
		"countryName":   "Germany",
	}

	err := validator.ValidateAndSanitize(data)
	assert.NoError(t, err)
}

func TestBankRequestValidator_MissingFields(t *testing.T) {
	validator := NewBankRequestValidator()

	// No address
	data := map[string]interface{}{
		"bankName":      "Test Bank",
		"countryISO2":   "DE",
		"swiftCode":     "DEUTDEFF",
		"codeType":      "BIC8",
		"isHeadquarter": false,
		"timeZone":      "Europe/Berlin",
	}
	err := validator.ValidateAndSanitize(data)
	assert.EqualError(t, err, "address is required")

	data["address"] = "Some Address"
	data["bankName"] = "   "
	err = validator.ValidateAndSanitize(data)
	assert.EqualError(t, err, "bankName must be a non-empty string")
}

func TestBankRequestValidator_InvalidCountryCode(t *testing.T) {
	validator := NewBankRequestValidator()

	data := map[string]interface{}{
		"address":       "123 Main St",
		"bankName":      "Bank",
		"countryISO2":   "XYZ",
		"swiftCode":     "DEUTXYZF",
		"codeType":      "BIC8",
		"isHeadquarter": false,
		"timeZone":      "Europe/Berlin",
		"countryName":   "Germany",
	}
	err := validator.ValidateAndSanitize(data)
	assert.Contains(t, err.Error(), "countryISO2 invalid")
}

func TestBankRequestValidator_InvalidSwiftCode(t *testing.T) {
	validator := NewBankRequestValidator()

	data := map[string]interface{}{
		"address":       "123 Main St",
		"bankName":      "Bank",
		"countryISO2":   "DE",
		"swiftCode":     "INVALID",
		"codeType":      "BIC8",
		"isHeadquarter": false,
		"timeZone":      "Europe/Berlin",
		"countryName":   "Germany",
	}
	err := validator.ValidateAndSanitize(data)
	assert.Contains(t, err.Error(), "swiftCode invalid")
}

func TestBankRequestValidator_SwiftCodeMismatch(t *testing.T) {
	validator := NewBankRequestValidator()

	data := map[string]interface{}{
		"address":       "123 Main St",
		"bankName":      "Bank",
		"countryISO2":   "FR",
		"swiftCode":     "DEUTDEFF",
		"codeType":      "BIC8",
		"isHeadquarter": false,
		"timeZone":      "Europe/Paris",
		"countryName":   "France",
	}
	err := validator.ValidateAndSanitize(data)
	assert.Contains(t, err.Error(), "swiftCode and countryISO2 mismatch")
}

func TestBankRequestValidator_HeadquarterMismatch(t *testing.T) {
	validator := NewBankRequestValidator()

	data := map[string]interface{}{
		"address":       "123 Main St",
		"bankName":      "Bank",
		"countryISO2":   "DE",
		"swiftCode":     "DEUTDEFF", // BIC8 for diversity
		"codeType":      "BIC8",
		"isHeadquarter": true,
		"timeZone":      "Europe/Berlin",
		"countryName":   "Germany",
	}
	err := validator.ValidateAndSanitize(data)
	assert.Contains(t, err.Error(), "isHeadquarter and swift code mismatch")
	data = map[string]interface{}{
		"address":       "123 Main St",
		"bankName":      "Bank",
		"countryISO2":   "DE",
		"swiftCode":     "DEUTDEFFXXX",
		"codeType":      "BIC8",
		"isHeadquarter": false,
		"timeZone":      "Europe/Berlin",
		"countryName":   "Germany",
	}
	err = validator.ValidateAndSanitize(data)
	assert.Contains(t, err.Error(), "isHeadquarter and swift code mismatch")
}
