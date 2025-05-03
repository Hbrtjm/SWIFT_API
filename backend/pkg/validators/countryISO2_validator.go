package validators

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// CountryCodeValidator validates country codes
type CountryISO2CodeValidator struct{}

func NewCountryISO2CodeValidator() *CountryISO2CodeValidator {
	return &CountryISO2CodeValidator{}
}

func (ccv *CountryISO2CodeValidator) Validate(value interface{}) error {
	countryISO2, ok := value.(string)
	if !ok {
		return errors.New("country code must be a string")
	}

	countryISO2 = strings.ToUpper(countryISO2)
	if len(countryISO2) < 2 {
		return fmt.Errorf("country code is too short: %d", len(countryISO2))
	}
	if len(countryISO2) > 2 {
		return fmt.Errorf("country code is too long: %d", len(countryISO2))
	}

	re := regexp.MustCompile("^[A-Z]{2}$")
	if !re.MatchString(countryISO2) {
		return fmt.Errorf("invalid country code format: %s", countryISO2)
	}

	return nil
}
