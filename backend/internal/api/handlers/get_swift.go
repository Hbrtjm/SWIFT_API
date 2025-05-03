package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// GetBySwiftCode handles GET request for a single SWIFT code
func (rh *RequestsHandler) GetBySwiftCode(w http.ResponseWriter, r *http.Request) {
	// Extract the SWIFT code from the URL params
	vars := mux.Vars(r)
	swiftCode := vars["swiftCode"]

	// Cannot be empty
	if swiftCode == "" {
		rh.logger.Error("SWIFT code is empty")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "SWIFT code is required"})
		return
	}

	rh.logger.Debug("Getting by SWIFT code: %s", swiftCode)

	response, err := rh.service.GetBySwiftCode(swiftCode)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		errResponse := map[string]string{"error": "Country code not found"}
		if IsAPIDebugActive() {
			errResponse["error"] = err.Error()
		}
		w.WriteHeader(http.StatusNotFound)
		rh.logger.Error("Error fetching SWIFT code: %v", err)
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
