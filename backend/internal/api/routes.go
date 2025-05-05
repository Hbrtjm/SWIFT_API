package api

import (
	"fmt"
	"net/http"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/api/handlers"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/api/middleware"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/service"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/util"
	"github.com/gorilla/mux"
)

// NewRouter creates and configures a new router with all API routes
func NewRouter(service *service.SwiftCodeService, logger *middleware.Logger) (router *mux.Router) {
	router = mux.NewRouter()

	// Create the request handler with the new logger interface
	swiftDatabaseResponseHandler := handlers.NewRequestsHandler(service, logger)

	// Use the new middleware with default configuration
	config := middleware.DefaultConfig()
	config.LogRequestBody = true
	config.LogResponseBody = true

	router.Use(middleware.LoggingMiddlewareWithConfig(logger, config))
	router.Use(middleware.ContentTypeMiddleware)

	version := util.GetEnvOrDefault("VERSION", "v1")

	api := router.PathPrefix(fmt.Sprintf("/%s", version)).Subrouter()

	// Define all API routes
	api.HandleFunc("/swift-codes/{swiftCode}", swiftDatabaseResponseHandler.GetBySwiftCode).Methods(http.MethodGet)
	api.HandleFunc("/swift-codes/country/{countryISO2code}", swiftDatabaseResponseHandler.GetBySwiftCodesByCountry).Methods(http.MethodGet)
	api.HandleFunc("/swift-codes", swiftDatabaseResponseHandler.PostBankEntry).Methods(http.MethodPost)
	api.HandleFunc("/swift-codes/{swiftCode}", swiftDatabaseResponseHandler.DeleteSwiftCode).Methods(http.MethodDelete)

	// Health check endpoint, used in testing
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Health check requested")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	return router
}
