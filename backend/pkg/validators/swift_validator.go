package validators

import (
	"errors"
	"fmt"
	"regexp"
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

	if len(swiftCode) != 8 && len(swiftCode) != 11 {
		return fmt.Errorf("invalid SWIFT code length: %d", len(swiftCode))
	}

	re := regexp.MustCompile(`^[A-Z]{4}[A-Z]{2}[A-Z0-9]{2}([A-Z0-9]{3})?$`)
	if !re.MatchString(swiftCode) {
		return fmt.Errorf("invalid SWIFT code format: %s", swiftCode)
	}

	return nil
}

func (sv *SwiftCodeValidator) ValidateWithCountryCode(value interface{}, countryISO2 string) error {
	err := sv.Validate(value)
	if err != nil {
		return err
	}

	swiftCode, ok := value.(string)
	if !ok {
		return fmt.Errorf("SWIFT code must be a string, got: %s", swiftCode)
	}

	// Check if countryISO2 matches SWIFT positions 5-6
	if swiftCode[4:6] != countryISO2 {
		return fmt.Errorf("SWIFT code country does not match provided country code: %s vs %s", swiftCode[4:6], countryISO2)
	}

	return nil
}

// Generally it is assumed that the the isHeadquarter field would be evaluated on our side
func (sv *SwiftCodeValidator) ValidateWithIsHeadquarter(code string, isHeadquarter bool) error {

	err := sv.Validate(code)

	if err != nil {
		return fmt.Errorf("error checking if bank is headquarter: %w", err)
	}
	if code[8:] == "XXX" && !isHeadquarter {
		return errors.New("the bank should be marked as a headquarter")
	}

	if code[8:] != "XXX" && isHeadquarter {
		return errors.New("the bank is not a headquarter")
	}

	return nil
}
