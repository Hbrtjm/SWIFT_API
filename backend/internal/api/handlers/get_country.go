package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Returns all SWIFT codes for a given country
func (rh *RequestsHandler) GetBySwiftCodesByCountry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryISO2 := vars["countryISO2code"]

	if countryISO2 == "" {
		rh.logger.Error("Country ISO2 code is empty")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Country ISO2 code cannot be empty"})
	}

	rh.logger.Info("Getting SWIFT codes for country: %s", countryISO2)

	response, err := rh.service.GetBySwiftCodesByCountry(countryISO2)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		errResponse := map[string]string{"message": "Country code not found"}
		if IsAPIDebugActive() {
			errResponse["message"] = err.Error()
		}
		w.WriteHeader(http.StatusNotFound)
		rh.logger.Error("Error fetching country code: %v", err)
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
