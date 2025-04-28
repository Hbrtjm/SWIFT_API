package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// GetBySwiftCodesByCountry returns all SWIFT codes for a given country
func (rh *RequestsHandler) GetBySwiftCodesByCountry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryCode := vars["countryISO2code"]

	rh.logger.Printf("Getting SWIFT codes for country: %s", countryCode)

	response, err := rh.service.GetBySwiftCodesByCountry(countryCode)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		rh.logger.Printf("Error fetching country code: %v", err)
		json.NewEncoder(w).Encode(map[string]string{"error": "Country code not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
