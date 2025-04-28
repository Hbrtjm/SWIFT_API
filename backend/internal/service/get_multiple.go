package service

import (
	"errors"
)

// GetMultipleSwiftCodes returns data for multiple SWIFT codes
func (s *SwiftCodeService) GetMultipleSwiftCodes(codes []string) ([]map[string]interface{}, error) {
	if len(codes) == 0 {
		return []map[string]interface{}{}, nil
	}

	result := make([]map[string]interface{}, 0, len(codes))
	for _, code := range codes {
		data, err := s.repo.FindBySwiftCode(code)
		if err == nil {
			result = append(result, data)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no valid SWIFT codes found")
	}

	return result, nil
}
