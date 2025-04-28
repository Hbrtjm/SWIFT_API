package validators

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type TimeZoneValidator struct{}

func (tzv *TimeZoneValidator) Validate(value interface{}, countryName string) error {
	timeZone, ok := value.(string)
	if !ok {
		return errors.New("timezone value must be a string")
	}

	timeZone = strings.ToUpper(timeZone)

	countryName = strings.ToLower(countryName)

	// TODO - Placeholder, I need to create a dictionary of country names to their capital names
	capitalName := strings.ToUpper(string(countryName[0])) + countryName[1:]

	re := regexp.MustCompile(`^(Europe|Asia|America)/` + capitalName + `$`)
	if !re.MatchString(timeZone) {
		return fmt.Errorf("invalid timezone format: %s", timeZone)
	}

	return nil
}
