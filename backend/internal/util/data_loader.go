package util

import (
	"github.com/Hbrtjm/SWIFT_API/backend/internal/api/middleware"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/repository"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/service"
)

// LoadInitialDataIfNeeded checks if database is empty and loads initial data if needed
func LoadInitialDataIfNeeded(service *service.SwiftCodeService, repo *repository.MongoRepository, filename string, logger *middleware.Logger) error {
	count, err := repo.Count()

	if err != nil {
		logger.Error("Error checking database: %v", err)
		return err
	}

	if count > 0 {
		logger.Warning("Database already contains %d SWIFT codes, inserting the contents from the file", count)
	}

	// logger.Info("Database is empty. Loading initial SWIFT data from %s", filename)
	err = service.LoadInitialData(filename)
	if err != nil {
		logger.Error("Error loading initial data: %v", err)
		return err
	}

	newRowCount, _ := repo.Count()
	logger.Info("Successfully loaded %d SWIFT codes into the database", newRowCount-count)

	return nil
}
