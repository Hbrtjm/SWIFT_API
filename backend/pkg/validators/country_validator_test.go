package validators

import (
	"testing"
)

func TestCountryValidator_ValidInput(t *testing.T) {
	validator := NewCountryValidator()

	input := map[string]interface{}{
		"codeType":    "BIC11",
		"timeZone":    "Europe/Berlin",
		"countryName": "Germany",
	}

	err := validator.ValidateAndSanitize(input)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestCountryValidator_MissingCodeType(t *testing.T) {
	validator := NewCountryValidator()

	input := map[string]interface{}{
		"timeZone":    "Europe/Berlin",
		"countryName": "Germany",
	}

	err := validator.ValidateAndSanitize(input)
	if err == nil {
		t.Error("Expected error for missing codeType, got nil")
	}
}

func TestCountryValidator_EmptyCountryName(t *testing.T) {
	validator := NewCountryValidator()

	input := map[string]interface{}{
		"codeType":    "BiC8",
		"timeZone":    "Europe/Berlin",
		"countryName": "   ",
	}

	err := validator.ValidateAndSanitize(input)
	if err == nil {
		t.Error("Expected error for empty countryName, got nil")
	}
}

func TestCountryValidator_InvalidCodeTypeFormat(t *testing.T) {
	validator := NewCountryValidator()

	input := map[string]interface{}{
		"codeType":    "INVALID123",
		"timeZone":    "Europe/Berlin",
		"countryName": "Germany",
	}

	err := validator.ValidateAndSanitize(input)
	if err == nil {
		t.Error("Expected error for invalid codeType format, got nil")
	}
}

func TestCountryValidator_InvalidTimeZone(t *testing.T) {
	validator := NewCountryValidator()

	input := map[string]interface{}{
		"codeType":    "Bic11",
		"timeZone":    "Europe/Warsaw21", // assumed invalid timezone
		"countryName": "Germany",
	}

	err := validator.ValidateAndSanitize(input)
	if err == nil {
		t.Error("Expected error for invalid timeZone, got nil")
	}
}

func TestCountryValidator_InjectionCharacters(t *testing.T) {
	validator := NewCountryValidator()

	input := map[string]interface{}{
		"codeType":    "BIc8",
		"timeZone":    "Europe/Berlin",
		"countryName": "Germ${}any",
	}

	err := validator.ValidateAndSanitize(input)
	if err == nil {
		t.Error("Expected error for injection characters, got nil")
	}
}
