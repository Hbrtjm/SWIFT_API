package validators

import (
	"fmt"
	"regexp"
	"strings"
)

func Sanitize(data map[string]interface{}) error {
	// No field should contain $, { or }, since that could lead to a MongoDB injection
	illegalCharacters := regexp.MustCompile(`(\$|\}|\{)`)
	var validationErrors []string

	for key, fieldInterface := range data {
		field, ok := fieldInterface.(string)
		if !ok {
			continue
		}
		if illegalCharacters.MatchString(field) {
			validationErrors = append(validationErrors, fmt.Sprintf("field %s contains illegal value: %s", key, field))
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("validation errors:\n%s", strings.Join(validationErrors, "\n"))
	}
	return nil
}
