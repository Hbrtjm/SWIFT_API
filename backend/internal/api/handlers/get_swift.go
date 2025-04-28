package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// GetBySwiftCode handles GET request for a single SWIFT code
func (rh *RequestsHandler) GetBySwiftCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	swiftCode := vars["swiftCode"]

	rh.logger.Printf("Getting by SWIFT code: %s", swiftCode)

	response, err := rh.service.GetBySwiftCode(swiftCode)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		rh.logger.Printf("Error fetching SWIFT code: %v", err)
		json.NewEncoder(w).Encode(map[string]string{"error": "SWIFT code not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
