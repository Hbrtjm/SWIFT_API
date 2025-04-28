package validators

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// CountryCodeValidator validates country codes
type CountryCodeValidator struct{}

func NewCountryCodeValidator() *CountryCodeValidator {
	return &CountryCodeValidator{}
}

func (ccv *CountryCodeValidator) Validate(value interface{}) error {
	countryCode, ok := value.(string)
	if !ok {
		return errors.New("country code value must be a string")
	}

	countryCode = strings.ToUpper(countryCode)
	if len(countryCode) != 2 {
		return fmt.Errorf("invalid country code length: %d", len(countryCode))
	}

	re := regexp.MustCompile("^[A-Z]{2}$")
	if !re.MatchString(countryCode) {
		return fmt.Errorf("invalid country code format: %s", countryCode)
	}

	return nil
}
