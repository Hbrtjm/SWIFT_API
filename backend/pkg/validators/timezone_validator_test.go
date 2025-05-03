package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeZoneValidator(t *testing.T) {
	validator := NewTimeZoneValidator()

	assert.NoError(t, validator.Validate("Europe/Warsaw", "Poland"))
	assert.NoError(t, validator.Validate("America/NewYork", "USA"))
	assert.NoError(t, validator.Validate("EUROPE/WARSAW", "Poland"))

	assert.EqualError(t, validator.Validate(123, "USA"), "timeZone value must be a string")

	assert.EqualError(t, validator.Validate("EuropeWarsaw", "Poland"), "invalid timeZone format: EUROPEWARSAW")
}
