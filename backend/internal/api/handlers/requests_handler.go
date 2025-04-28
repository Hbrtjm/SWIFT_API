package handlers

import (
	"log"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/service"
)

// RequestsHandler handles HTTP requests related to SWIFT codes
type RequestsHandler struct {
	service *service.SwiftCodeService
	logger  *log.Logger
}

// NewRequestsHandler creates a new RequestsHandler
func NewRequestsHandler(service *service.SwiftCodeService, logger *log.Logger) *RequestsHandler {
	return &RequestsHandler{
		service: service,
		logger:  logger,
	}
}
