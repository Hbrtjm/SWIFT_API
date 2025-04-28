// This is an unused validator, it's for more general purposes abstracting away particular validators, and creating a lookup dictionary for validations

package validator

import (
	"github.com/Hbrtjm/SWIFT_API/backend/pkg/validators"
)

type GeneralValidator interface {
	Validate(value interface{}) error
}

func init() {
	Register("swift", &validators.SwiftCodeValidator{})
	Register("countryCode", &validators.CountryCodeValidator{})

}

func Vlidate(name string, value interface{}) error {
	validator, err := GetValidator(name)
	if err != nil {
		return err
	}
	return validator.Validate(value)
}
