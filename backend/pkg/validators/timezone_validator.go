package validators

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type TimeZoneValidator struct{}

func NewTimeZoneValidator() *TimeZoneValidator {
	return &TimeZoneValidator{}
}

func (tzv *TimeZoneValidator) Validate(value interface{}, countryName string) error {
	timeZone, ok := value.(string)
	if !ok {
		return errors.New("timeZone value must be a string")
	}

	timeZone = strings.ToUpper(timeZone)

	capitalName := `[A-Z]{1}[a-zA-Z]+`

	re := regexp.MustCompile(`^[A-Za-z]+/` + capitalName + `$`)
	if !re.MatchString(timeZone) {
		return fmt.Errorf("invalid timeZone format: %s", timeZone)
	}

	return nil
}
