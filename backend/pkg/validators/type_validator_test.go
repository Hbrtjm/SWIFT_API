package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodeTypeValidator(t *testing.T) {
	validator := NewCodeTypeValidator()

	// Correct format, but in our current data we don't have any BIC-s 8
	assert.NoError(t, validator.Validate("BIC8"))
	assert.NoError(t, validator.Validate("bic11"))

	assert.EqualError(t, validator.Validate(123), "code type value must be a string")
	assert.EqualError(t, validator.Validate("XYZ"), "invalid code type format: XYZ")
	assert.EqualError(t, validator.Validate("BIC12"), "invalid code type format: BIC12")
}
