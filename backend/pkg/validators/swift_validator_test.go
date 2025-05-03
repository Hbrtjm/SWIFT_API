package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwiftCodeValidator(t *testing.T) {
	validator := NewSwiftCodeValidator()

	assert.NoError(t, validator.Validate("DEUTDEFF"))
	assert.NoError(t, validator.Validate("DEUTDEFF500"))

	assert.EqualError(t, validator.Validate(123), "swift code value must be a string")

	assert.EqualError(t, validator.Validate("DEUT"), "invalid SWIFT code length: 4")

	assert.EqualError(t, validator.Validate("abcd1234"), "invalid SWIFT code format: abcd1234")
}

func TestSwiftCodeValidator_CountryCode(t *testing.T) {

	validator := NewSwiftCodeValidator()

	err := validator.ValidateWithCountryCode("DEUTUS33", "DE")
	assert.EqualError(t, err, "SWIFT code country does not match provided country code: US vs DE")

	assert.NoError(t, validator.ValidateWithCountryCode("DEUTDEFF", "DE"))
}

func TestSwiftCodeValidator_HeadquarterCheck(t *testing.T) {
	validator := NewSwiftCodeValidator()

	assert.NoError(t, validator.ValidateWithIsHeadquarter("DEUTDEFFXXX", true))
	assert.EqualError(t, validator.ValidateWithIsHeadquarter("DEUTDEFFXXX", false), "the bank should be marked as a headquarter")
	assert.EqualError(t, validator.ValidateWithIsHeadquarter("DEUTDEFF500", true), "the bank is not a headquarter")
}
