package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Handles DELETE request for a SWIFT code
func (rh *RequestsHandler) DeleteSwiftCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	swiftCode := vars["swiftCode"]

	err := rh.service.DeleteSwiftCode(swiftCode)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		errResponse := map[string]string{"error": "Failed to delete SWIFT code"}
		if IsAPIDebugActive() {
			errResponse["error"] = err.Error()
		}
		w.WriteHeader(http.StatusNotFound)
		rh.logger.Error("Failed to delete SWIFT code:  %v", err)
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "SWIFT code deleted successfully"})
}
