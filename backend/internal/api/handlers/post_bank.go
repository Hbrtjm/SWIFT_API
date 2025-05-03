package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/util"
)

// PostBankEntry handles POST request to create a new SWIFT code
func (rh *RequestsHandler) PostBankEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Define a struct that matches the expected JSON format
	swiftCodeRequest := util.NewPostRequest()

	// Decode the request body into the struct
	err := json.NewDecoder(r.Body).Decode(&swiftCodeRequest)
	if err != nil {
		rh.logger.Error("Invalid request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request body"})
		return
	}

	bankData := map[string]interface{}{
		"swiftCode":     strings.ToUpper(swiftCodeRequest.SwiftCode),
		"countryISO2":   strings.ToUpper(swiftCodeRequest.CountryISO2),
		"countryName":   strings.ToUpper(swiftCodeRequest.CountryName),
		"bankName":      swiftCodeRequest.BankName,
		"address":       swiftCodeRequest.Address,
		"townName":      swiftCodeRequest.TownName,
		"timeZone":      swiftCodeRequest.TimeZone,
		"isHeadquarter": swiftCodeRequest.IsHeadquarter,
		"codeType":      swiftCodeRequest.CodeType,
	}

	// Pass bank data to service layer for processing
	err = rh.service.PostBankData(bankData)

	if err != nil {

		rh.logger.Error("Error creating bank: %v", err)

		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "validation error") ||
			strings.Contains(err.Error(), "already exists") {
			statusCode = http.StatusBadRequest
		}

		errResponse := map[string]string{"message": "Country code not found"}
		if IsAPIDebugActive() {
			errResponse["message"] = err.Error()
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	// Success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Bank with SWIFT code created successfully",
	})
}
