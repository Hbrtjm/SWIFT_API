package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	// TODO - For later if I want to add frontend
	// "github.com/rs/cors"
	"github.com/Hbrtjm/SWIFT_API/internal/api/handlers"
	"github.com/Hbrtjm/SWIFT_API/internal/api/middleware"
	"github.com/Hbrtjm/SWIFT_API/internal/service"
)

func NewRouter(service *service.SwiftCodeService, logger *log.Logger) (router *mux.Router) {
	router = mux.NewRouter()

	swiftParser := handlers.NewSwiftParserHandler(service, logger)

	router.Use(middleware.LogginMiddleware(logger))
	router.Use(middleware.ContentTypeMiddleware)

	api := router.PathPrefix(fmt.Sprintf("/%s", os.Getenv("VERSION"))).Subrouter()

	api.HandleFunc("/swift-codes/{swiftCode}", swiftParser.GetSwiftCode).Methods(http.MethodGet)
	api.HandleFunc("/swift-codes/country/{countryISO2code}", swiftParser.GetSwiftCodesByCountry).Methods(http.MethodGet)
	api.HandleFunc("/swift-codes", swiftParser.CreateBankEntry).Methods(http.MethodPost)
	api.HandleFunc("/swift-codes/{swiftCode}", swiftParser.DeleteSwiftCode).Methods(http.MethodDelete)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	return router
}
