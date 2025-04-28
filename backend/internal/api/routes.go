package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/api/handlers"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/api/middleware"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/service"
	"github.com/gorilla/mux"
)

// NewRouter creates and configures a new router with all API routes
func NewRouter(service *service.SwiftCodeService, logger *log.Logger) (router *mux.Router) {
	router = mux.NewRouter()

	swiftDatabaseResponseHandler := handlers.NewRequestsHandler(service, logger)

	router.Use(middleware.LogginMiddleware(logger))
	router.Use(middleware.ContentTypeMiddleware)

	// API routes with versioning
	api := router.PathPrefix(fmt.Sprintf("/%s", os.Getenv("VERSION"))).Subrouter()

	// Define all API routes
	api.HandleFunc("/swift-codes/{swiftCode}", swiftDatabaseResponseHandler.GetBySwiftCode).Methods(http.MethodGet)
	api.HandleFunc("/swift-codes/country/{countryISO2code}", swiftDatabaseResponseHandler.GetBySwiftCodesByCountry).Methods(http.MethodGet)
	api.HandleFunc("/swift-codes", swiftDatabaseResponseHandler.PostBankEntry).Methods(http.MethodPost)
	api.HandleFunc("/swift-codes/{swiftCode}", swiftDatabaseResponseHandler.DeleteSwiftCode).Methods(http.MethodDelete)

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	return router
}
