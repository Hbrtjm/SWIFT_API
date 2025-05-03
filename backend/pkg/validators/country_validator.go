package validators

import (
	"fmt"
	"strings"
)

type CountryValidator struct {
	codeTypeValidator        *CodeTypeValidator
	timeZoneValidator        *TimeZoneValidator
	countryISO2CodeValidator *CountryISO2CodeValidator
}

func NewCountryValidator() *CountryValidator {
	return &CountryValidator{
		codeTypeValidator:        NewCodeTypeValidator(),
		timeZoneValidator:        NewTimeZoneValidator(),
		countryISO2CodeValidator: NewCountryISO2CodeValidator(),
	}
}

func (cv *CountryValidator) ValidateAndSanitize(data map[string]interface{}) error {
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

	err := Sanitize(data)
	if err != nil {
		return err
	}

	// Required fields
	// I don't have valid ideas for address verification and I don't want to make it to complex
	if _, err := getString("countryName"); err != nil {
		return err
	}

	codeType, err := getString("codeType")
	if err != nil {
		return err
	}
	timeZone, err := getString("timeZone")
	if err != nil {
		return err
	}
	countryName, _ := getString("countryName")

	if err := cv.codeTypeValidator.Validate(codeType); err != nil {
		return fmt.Errorf("codeType invalid: %w", err)
	}

	if err := cv.timeZoneValidator.Validate(timeZone, countryName); err != nil {
		return fmt.Errorf("timeZone invalid: %w", err)
	}

	return nil
}
