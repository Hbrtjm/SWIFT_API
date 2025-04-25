package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Hbrtjm/SWIFT_API/internal/service"
	"github.com/gorilla/mux"
)

// SwiftParserHandler handles HTTP requests related to SWIFT codes
type SwiftParserHandler struct {
	service *service.SwiftCodeService
	logger  *log.Logger
}

func NewSwiftParserHandler(service *service.SwiftCodeService, logger *log.Logger) *SwiftParserHandler {
	return &SwiftParserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *SwiftParserHandler) GetSwiftCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	swiftCode := vars["swiftCode"]

	h.logger.Printf("Getting SWIFT code: %s", swiftCode)

	response, err := h.service.GetSwiftCode(swiftCode)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.logger.Printf("Error fetching SWIFT code: %v", err)
		json.NewEncoder(w).Encode(map[string]string{"error": "SWIFT code not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetSwiftCodesByCountry returns all SWIFT codes for a given country
func (h *SwiftParserHandler) GetSwiftCodesByCountry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryCode := vars["countryISO2code"]

	h.logger.Printf("Getting SWIFT codes for country: %s", countryCode)

	response, err := h.service.GetSwiftCodesByCountry(countryCode)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.logger.Printf("Error fetching country code: %v", err)
		json.NewEncoder(w).Encode(map[string]string{"error": "Country code not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetSwiftCodes handles POST request to fetch multiple SWIFT codes at once
func (h *SwiftParserHandler) GetSwiftCodes(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		SwiftCodes []string `json:"swiftCodes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Printf("Getting multiple SWIFT codes: %v", requestBody.SwiftCodes)

	// Dummy response with requested codes
	response := make([]map[string]interface{}, 0, len(requestBody.SwiftCodes))
	for _, code := range requestBody.SwiftCodes {
		response = append(response, map[string]interface{}{
			"swiftCode":    code,
			"bankName":     "Bank for " + code,
			"countryCode":  "US",
			"countryName":  "UNITED STATES",
			"city":         "New York",
			"branchCode":   code[8:],
			"isHeadOffice": code[8:] == "XXX",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateSwiftCode handles the creation of a new SWIFT code
func (h *SwiftParserHandler) DeleteSwiftCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	swiftCode := vars["swiftCode"]

	err := h.service.DeleteSwiftCode(swiftCode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Printf("Error deleting SWIFT code: %v", err)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete SWIFT code"})
		return
	}

	// Dummy response for delete operation
	response := map[string]interface{}{
		"message": "SWIFT code deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *SwiftParserHandler) CreateBankEntry(w http.ResponseWriter, r *http.Request) {
	var bankData map[string]interface{}
	// TODO - For now we are doing this simply, but I have to elimitate IsHeadOffice and BranchCode from the request body, as they should be implied by the rest of the data
	err := json.NewDecoder(r.Body).Decode(&bankData)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Printf("Invalid request body: %v", err)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	err = h.service.CreateSwiftCode(bankData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Printf("Error creating SWIFT code: %v", err)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create SWIFT code"})
		return
	}

	response := map[string]interface{}{
		"message": "SWIFT code created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
