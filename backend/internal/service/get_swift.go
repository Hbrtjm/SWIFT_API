package service

import (
	"errors"
)

type SwiftCodeResponse struct {
	Address       string                   `json:"address"`
	BankName      string                   `json:"bankName"`
	CountryISO2   string                   `json:"countryISO2"`
	CountryName   string                   `json:"countryName"`
	IsHeadquarter bool                     `json:"isHeadquarter"`
	SwiftCode     string                   `json:"swiftCode"`
	Branches      []map[string]interface{} `json:"branches,omitempty"`
}

func (s *SwiftCodeService) GetBySwiftCode(code string) (*SwiftCodeResponse, error) {
	value, err := s.repo.FindBySwiftCode(code)
	if err != nil || value == nil {
		return nil, errors.New("no bank found with the given SWIFT code")
	}

	// Check if this is a headquarter
	isHeadquarter, _ := value["isHeadquarter"].(bool)
	if isHeadquarter {
		// Get the branch code
		branchCode, ok := value["branchCode"].(string)
		if ok && branchCode != "" {
			// Find all branches with the same branch code
			branchesList, err := s.repo.FindByBranchCode(branchCode)
			if err != nil {
				s.logger.Printf("Error finding branches: %v", err)
				// Return the headquarter info even if there was an error finding branches
				return mapToSwiftCodeResponse(value), nil
			}

			// Filter out the headquarter itself from the branch list
			swiftCode, _ := value["swiftCode"].(string)
			filteredBranches := make([]map[string]interface{}, 0)
			for _, branch := range branchesList {
				branchSwiftCode, ok := branch["swiftCode"].(string)
				if ok && branchSwiftCode != swiftCode {
					filteredBranches = append(filteredBranches, branch)
				}
			}

			// Create a new SwiftCodeResponse with the branches included
			response := mapToSwiftCodeResponse(value)
			if len(filteredBranches) > 0 {
				// Add branches to the response only if any are found
				response.Branches = filteredBranches
			}
			return response, nil
		}
	}

	// Return the original value if it's not a headquarter or has no branches
	return mapToSwiftCodeResponse(value), nil
}

// Helper function to map a map[string]interface{} to SwiftCodeResponse
func mapToSwiftCodeResponse(value map[string]interface{}) *SwiftCodeResponse {
	address, _ := value["address"].(string)
	bankName, _ := value["bankName"].(string)
	countryISO2, _ := value["countryISO2"].(string)
	countryName, _ := value["countryName"].(string)
	isHeadquarter, _ := value["isHeadquarter"].(bool)
	swiftCode, _ := value["swiftCode"].(string)

	return &SwiftCodeResponse{
		Address:       address,
		BankName:      bankName,
		CountryISO2:   countryISO2,
		CountryName:   countryName,
		IsHeadquarter: isHeadquarter,
		SwiftCode:     swiftCode,
	}
}
