package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/util"
)

// PostBankEntry handles POST request to create a new SWIFT code
func (rh *RequestsHandler) PostBankEntry(w http.ResponseWriter, r *http.Request) {
	// Define a struct that matches the expected JSON format
	swiftCodeRequest := util.NewPostRequest()

	// Decode the request body into the struct
	err := json.NewDecoder(r.Body).Decode(&swiftCodeRequest)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		rh.logger.Printf("Invalid request body: %v", err)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate the SWIFT code
	if len(swiftCodeRequest.SwiftCode) < 8 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "SWIFT code must be at least 8 characters"})
		return
	}

	// Convert to map to work with existing service
	bankData := map[string]interface{}{
		"address":       swiftCodeRequest.Address,
		"bankName":      swiftCodeRequest.BankName,
		"countryCode":   strings.ToUpper(swiftCodeRequest.CountryISO2),
		"countryName":   strings.ToUpper(swiftCodeRequest.CountryName),
		"isHeadquarter": swiftCodeRequest.SwiftCode[8:] == "XXX",
		"swiftCode":     swiftCodeRequest.SwiftCode,
		"branchCode":    swiftCodeRequest.SwiftCode[:8],
	}

	// TODO - Maybe I could could work on the same address of bank data
	err = rh.service.PostBankData(bankData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		rh.logger.Printf("Error creating SWIFT code: %v", err)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create SWIFT code"})
		return
	}

	// Set the content type and status before writing the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "SWIFT code created successfully",
	})
}
