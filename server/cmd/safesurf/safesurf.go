package main

import (
	"log"
	"os"

	"github.com/abhizaik/SafeSurf/internal/handler"
	"github.com/abhizaik/SafeSurf/internal/service/rank"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (non-fatal if missing)
	if err := godotenv.Load("/app/.env"); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	}

	r := handler.SetupRouter()

	err := rank.LoadDomainRanks()
	if err != nil {
		log.Fatal(err)
	}

	// Get port from environment variable, default to 8080
	port := getEnv("PORT", "8080")
	addr := ":" + port

	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
