package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hbrtjm/SWIFT_API/internal/api"
	"github.com/Hbrtjm/SWIFT_API/internal/db/repository"
	"github.com/Hbrtjm/SWIFT_API/internal/parser"
	"github.com/Hbrtjm/SWIFT_API/internal/service"
	"github.com/Hbrtjm/SWIFT_API/internal/util"
)

func main() {

	logger := log.New(os.Stdout, "SWIFT-API: ", log.LstdFlags)

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	debugMode := os.Getenv("DEBUG_MODE")

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "swiftcodes"
	}
	// We don't need that message in production, but it's useful for debugging
	if debugMode == "true" {
		logger.Printf("Connecting to MongoDB at %s", mongoURI)
	}

	repo, err := repository.NewMongoRepository(mongoURI, dbName)
	logger.Printf("Connected to MongoDB at %s %s", repo.Collection().Database().Name(), err)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer repo.Close()
	swiftFileParser := parser.NewSwiftCodeParser()
	swiftService := service.NewSwiftCodeService(repo, swiftFileParser, logger)

	// // Load initial data if needed
	if os.Getenv("LOAD_INITIAL_DATA") == "true" {
		filename := os.Getenv("SWIFT_DATA_FILE")
		if filename == "" {
			filename = "configs/swift_data.csv" // Default file path
		}

		err := util.LoadInitialDataIfNeeded(swiftService, repo, filename, logger)
		if err != nil {
			logger.Printf("Error with initial data process: %v", err)
		}

		// Create database indices for better performance
		err = util.CreateIndices(repo, logger)
		if err != nil {
			logger.Printf("Error creating database indices: %v", err)
		}
	}

	// ... and initialize the router with service
	router := api.NewRouter(swiftService, logger)

	// Configure server and run the server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		logger.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	// Create a deadline to wait for the server to shut down
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server exited properly")
}
