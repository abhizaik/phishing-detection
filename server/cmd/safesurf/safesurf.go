package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abhizaik/SafeSurf/internal/handler"
	"github.com/abhizaik/SafeSurf/internal/service/rank"
	"github.com/abhizaik/SafeSurf/internal/service/screenshot"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (non-fatal if missing)
	if err := godotenv.Load("/app/.env"); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	}

	// Initialize screenshot service (shared browser allocator)
	_, err := screenshot.GetService()
	if err != nil {
		log.Printf("Warning: Failed to initialize screenshot service: %v. Screenshot functionality may be unavailable.", err)
	} else {
		log.Println("Screenshot service initialized successfully")
		// Ensure cleanup on shutdown
		defer func() {
			service, _ := screenshot.GetService()
			if service != nil {
				service.Close()
			}
		}()
	}

	r := handler.SetupRouter()

	err = rank.LoadDomainRanks()
	if err != nil {
		log.Fatal(err)
	}

	// Get port from environment variable, default to 8080
	port := getEnv("PORT", "8080")
	addr := ":" + port

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on %s", addr)
		if err := r.Run(addr); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutting down server...")
	
	// Cleanup screenshot service
	service, _ := screenshot.GetService()
	if service != nil {
		service.Close()
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
