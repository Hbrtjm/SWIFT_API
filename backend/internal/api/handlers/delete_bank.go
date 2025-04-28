package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// DeleteSwiftCode handles DELETE request for a SWIFT code
func (rh *RequestsHandler) DeleteSwiftCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	swiftCode := vars["swiftCode"]

	err := rh.service.DeleteSwiftCode(swiftCode)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rh.logger.Printf("Error deleting SWIFT code: %v", err)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete SWIFT code"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "SWIFT code deleted successfully",
	})
}
