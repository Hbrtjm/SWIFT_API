package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/api"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/api/middleware"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/repository"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/parser"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/service"
	"github.com/Hbrtjm/SWIFT_API/backend/internal/util"
)

func main() {

	loggerPrefix := util.GetEnvOrDefault("LOGGER_PREFIX", "swiftcodes")

	logFileName := time.Now().Format("2007_04_10-15_04_05") + "-" + loggerPrefix + ".log"
	logger, err := middleware.FileDefaultLogger("./logs", logFileName, loggerPrefix)
	if err != nil {
		fmt.Printf("Failed to create a logger: %v\n", err) // Added newline
		logger = middleware.NewDefaultLogger(loggerPrefix)
	}

	speedupMode := util.GetEnvOrDefault("SPEEDUP_MODE", "false")
	if value, err := strconv.ParseBool(speedupMode); err == nil && value {
		logger = middleware.NewNoLogger()
	}

	// Connect to MongoDB
	mongoURI := util.GetEnvOrDefault("MONGO_URI", "mongodb://localhost:27017")

	dbName := util.GetEnvOrDefault("DB_NAME", "swiftcodes")

	banksCollectionName := util.GetEnvOrDefault("BANKS_COLLECTION_NAME", "banks")

	countriesCollectionName := util.GetEnvOrDefault("COUNTRIES_COLLECTION_NAME", "countries")

	// Debug information about MongoDB connection
	logger.Debug("Connecting to MongoDB at %s", mongoURI)

	repo, err := repository.NewMongoRepository(mongoURI, dbName, banksCollectionName, countriesCollectionName)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}

	logger.Info("Connected to MongoDB at %s", repo.BanksCollection().Database().Name())
	defer repo.CloseConnection()

	swiftFileParser := parser.NewSwiftFileParser()

	// Update the service to use our new logger
	swiftService := service.NewSwiftCodeService(repo, swiftFileParser, logger)

	// Load initial data if needed
	if util.GetEnvOrDefault("LOAD_INITIAL_DATA", "false") == "true" {
		filename := util.GetEnvOrDefault("SWIFT_DATA_FILE", "configs/swift_data.csv")

		err := util.LoadInitialDataIfNeeded(swiftService, repo, filename, logger)
		if err != nil {
			logger.Error("Error with initial data process: %v", err)
		}

		// Create database indices for better performance
		err = repo.CreateIndices(logger)
		if err != nil {
			logger.Error("Error creating database indices: %v", err)
		}
	}

	// Initialize the router with service and add our middleware
	router := api.NewRouter(swiftService, logger)

	// Configure and apply middleware
	middlewareConfig := middleware.DefaultConfig()
	middlewareConfig.LogRequestBody = true
	middlewareConfig.LogResponseBody = true

	// Configure server and run the server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline to wait for the server to shut down
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited properly")
}
