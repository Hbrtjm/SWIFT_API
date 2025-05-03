package handlers

import (
	"github.com/Hbrtjm/SWIFT_API/backend/internal/api/middleware"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/service"
)

// RequestsHandler handles HTTP requests related to SWIFT codes
type RequestsHandler struct {
	service *service.SwiftCodeService
	logger  *middleware.Logger
}

// NewRequestsHandler creates a new RequestsHandler
func NewRequestsHandler(service *service.SwiftCodeService, logger *middleware.Logger) *RequestsHandler {
	return &RequestsHandler{
		service: service,
		logger:  logger,
	}
}
