package validators

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type CodeTypeValidator struct{}

func NewCodeTypeValidator() *CodeTypeValidator {
	return &CodeTypeValidator{}
}

func (ctv *CodeTypeValidator) Validate(value interface{}) error {
	codeType, ok := value.(string)
	if !ok {
		return errors.New("code type value must be a string")
	}

	codeType = strings.ToUpper(codeType)

	re := regexp.MustCompile("^BIC(8|11)$") // In our database it's only BIC11, but we can handle BIC8 as well
	if !re.MatchString(codeType) {
		return fmt.Errorf("invalid code type format: %s", codeType)
	}

	return nil
}
