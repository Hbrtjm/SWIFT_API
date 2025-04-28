package validators

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type SwiftCodeValidator struct{}

func NewSwiftCodeValidator() *SwiftCodeValidator {
	return &SwiftCodeValidator{}
}

func (sv *SwiftCodeValidator) Validate(value interface{}) error {
	swiftCode, ok := value.(string)
	if !ok {
		return errors.New("swift code value must be a string")
	}

	swiftCode = strings.ToUpper(swiftCode)

	if len(swiftCode) != 8 && len(swiftCode) != 11 {
		return fmt.Errorf("invalid SWIFT code length: %d", len(swiftCode))
	}

	re := regexp.MustCompile(`^[A-Z]{4}[A-Z0-9]{2}[A-Z0-9]{2}([A-Z0-9]{3})?$`)
	if !re.MatchString(swiftCode) {
		return fmt.Errorf("invalid SWIFT code format: %s", swiftCode)
	}

	return nil
}

func (sv *SwiftCodeValidator) ValidateWithCountryCode(value interface{}, countryCode string) error {
	err := sv.Validate(value)
	if err != nil {
		return err
	}

	swiftCode := strings.ToUpper(value.(string))
	countryCode = strings.ToUpper(countryCode)

	// Check if countryCode matches SWIFT positions 5-6
	if swiftCode[4:6] != countryCode {
		return fmt.Errorf("SWIFT code country does not match provided country code: %s vs %s", swiftCode[4:6], countryCode)
	}

	return nil
}
