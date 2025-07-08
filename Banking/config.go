package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	PostgresURL string
	JWTKey      string
}

// LoadConfig loads environment variables and returns a Config struct
func LoadConfig() (*Config, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
		// Continue with environment variables from the system if .env file is not found
	}

	postgresURL := os.Getenv("Postgres_URL")
	if postgresURL == "" {
		log.Println("Warning: Postgres_URL is not set in environment variables")
	}

	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		log.Println("Warning: JWT_KEY is not set in environment variables")
	}

	return &Config{
		PostgresURL: postgresURL,
		JWTKey:      jwtKey,
	}, nil
}
