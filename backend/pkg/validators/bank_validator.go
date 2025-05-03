package validators

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type BankRequestValidator struct {
	swiftValidator    *SwiftCodeValidator
	countryValidator  *CountryISO2CodeValidator
	timeZoneValidator *TimeZoneValidator
	// codeTypeValidator *CodeTypeValidator
}

func NewBankRequestValidator() *BankRequestValidator {
	return &BankRequestValidator{
		swiftValidator:    NewSwiftCodeValidator(),
		countryValidator:  NewCountryISO2CodeValidator(),
		timeZoneValidator: NewTimeZoneValidator(),
		// codeTypeValidator: NewCodeTypeValidator(),
	}
}

func (v *BankRequestValidator) ValidateAndSanitize(data map[string]interface{}) error {
	// Helper function to get string value from map
	getString := func(key string) (string, error) {
		val, ok := data[key]
		if !ok {
			return "", fmt.Errorf("%s is required", key)
		}
		str, ok := val.(string)
		if !ok || strings.TrimSpace(str) == "" {
			return "", fmt.Errorf("%s must be a non-empty string", key)
		}
		return str, nil
	}

	// No field should contain $, { or }, since that could lead to a MongoDB injection
	illegalCharacters := regexp.MustCompile(`(\$|\}|\{})`)
	var validationErrors []string

	for key, fieldInterface := range data {
		field, ok := fieldInterface.(string)
		if !ok {
			continue
		}
		if illegalCharacters.MatchString(field) {
			validationErrors = append(validationErrors, strings.TrimSpace(fmt.Sprintf("field %s contains illegal value: %s", key, field)))
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("validation errors:\n%s", strings.Join(validationErrors, "\n"))
	}

	// Required fields
	// I don't have valid ideas for address verification and I don't want to make it to complex
	if _, err := getString("address"); err != nil {
		return err
	}
	// Looking through the data
	if _, err := getString("bankName"); err != nil {
		return err
	}
	countryISO2, err := getString("countryISO2")
	if err != nil {
		return err
	}
	swiftCode, err := getString("swiftCode")
	if err != nil {
		return err
	}

	// We can ignore any errors, if it's not present it will be infered from the SWFIT code itself
	isHeadquarter := getBool(data, "isHeadquarter")

	if err := v.countryValidator.Validate(countryISO2); err != nil {
		return fmt.Errorf("countryISO2 invalid: %w", err)
	}

	if err := v.swiftValidator.Validate(swiftCode); err != nil {
		return fmt.Errorf("swiftCode invalid: %w", err)
	}
	if err := v.swiftValidator.ValidateWithCountryCode(swiftCode, countryISO2); err != nil {
		return fmt.Errorf("swiftCode and countryISO2 mismatch: %w", err)
	}

	if err := v.swiftValidator.ValidateWithIsHeadquarter(swiftCode, isHeadquarter); err != nil {
		return fmt.Errorf("isHeadquarter and swift code mismatch: %w", err)
	}

	return nil
}

// Helper function to safely extract boolean values from map
func getBool(data map[string]interface{}, key string) bool {
	if value, exists := data[key]; exists {
		if boolValue, ok := value.(bool); ok {
			return boolValue
		}

		if strValue, ok := value.(string); ok {
			if boolValue, err := strconv.ParseBool(strValue); err == nil {
				return boolValue
			}
		}

		if numValue, ok := value.(float64); ok {
			return numValue != 0
		}
	}
	return false
}
