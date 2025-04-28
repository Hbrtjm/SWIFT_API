package service

import (
	"encoding/json"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/models"
)

// PostBankData creates a new bank entry in the database
func (s *SwiftCodeService) PostBankData(bankData map[string]interface{}) error {
	bankJson, err := json.Marshal(bankData)
	if err != nil {
		return err
	}

	var bank models.Bank
	err = json.Unmarshal(bankJson, &bank)
	if err != nil {
		return err
	}

	// Extract the branch code (first 8 characters of SWIFT code)
	if len(bank.SwiftCode) >= 8 {
		bank.BranchCode = bank.SwiftCode[:8]
	}

	return s.repo.Insert(bank)
}
