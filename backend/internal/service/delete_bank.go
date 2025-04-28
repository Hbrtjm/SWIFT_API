package service

import (
	"fmt"

	"github.com/Hbrtjm/SWIFT_API/backend/pkg/validators"
)

// DeleteSwiftCode deletes a SWIFT code from the database
func (s *SwiftCodeService) DeleteSwiftCode(code string) error {
	swiftValidator := validators.NewSwiftCodeValidator()
	swiftValidator.Validate(code)
	
	if err := swiftValidator.Validate(code); err != nil {
		return fmt.Errorf("invalid SWIFT code %s: %v", code, err)
	}

	return s.repo.Delete(code)
}
