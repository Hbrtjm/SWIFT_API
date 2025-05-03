package service

import (
	"errors"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/models"
)

func (s *SwiftCodeService) GetBySwiftCode(code string) (*SwiftCodeResponse, error) {
	bank, err := s.repo.FindBySwiftCode(code)
	emptyBank := models.Bank{}
	if err != nil || bank == emptyBank {
		return nil, errors.New("no bank found with the given SWIFT code")
	}

	// Get country name for the response
	countryName, err := s.repo.LookupCountryName(bank.CountryISO2)
	if err != nil {
		s.logger.Error("Error looking up country name: %v", err)
		countryName = "" // Continue even if country name lookup fails
	}

	if bank.IsHeadquarter {
		if bank.BranchCode != "" {
			branches, err := s.repo.FindByBranchCode(bank.BranchCode)
			if err != nil {
				s.logger.Error("Error finding branches: %v", err)
				// Return the headquarter info even if there was an error finding branches
				return bankToResponse(&bank, countryName), nil
			}

			// Filter out the headquarter itself from the branch list
			filteredBranches := make([]map[string]interface{}, 0)
			for _, branch := range branches {
				if branch["swiftCode"] != bank.SwiftCode {
					filteredBranches = append(filteredBranches, mapBranchValues(branch))
				}
			}

			// Create a new SwiftCodeResponse with the branches included
			response := bankToResponse(&bank, countryName)
			if len(filteredBranches) > 0 {
				response.Branches = filteredBranches
			}
			return response, nil
		}
	}

	// Return the original value if it's not a headquarter or has no branches
	return bankToResponse(&bank, countryName), nil
}
