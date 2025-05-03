package validators

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	validator := NewCountryISO2CodeValidator()

	correctValue := "PL"
	assert.Equal(t, nil, validator.Validate(correctValue))
	assert.NoError(t, validator.Validate(correctValue))
	assert.Equal(t, errors.New("country code must be a string"), validator.Validate(1))

	incorrectValue := "ABC"
	assert.Equal(t, fmt.Errorf("country code is too long: %d", len(incorrectValue)), validator.Validate(incorrectValue))
	incorrectValue = "A"
	assert.Equal(t, fmt.Errorf("country code is too short: %d", len(incorrectValue)), validator.Validate(incorrectValue))
}
