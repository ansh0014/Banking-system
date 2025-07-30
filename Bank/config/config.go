package config

import (
    "fmt"
    "os"
    // "path/filepath"

    "github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
    PostgresURL string
    JWTKey      string
}

// LoadConfig loads environment variables and returns a Config struct
func LoadConfig() (*Config, error) {
    // Try to load from project root first
    err := godotenv.Load()
    if err != nil {
        // If not found, try loading from current directory
        if err := godotenv.Load(".env"); err != nil {
            return nil, fmt.Errorf("error loading .env file: %v", err)
        }
    }

    postgresURL := os.Getenv("POSTGRES_URL")
    if postgresURL == "" {
        return nil, fmt.Errorf("POSTGRES_URL is not set in environment variables")
    }

    jwtKey := os.Getenv("JWT_KEY")
    if jwtKey == "" {
        return nil, fmt.Errorf("JWT_KEY is not set in environment variables")
    }

    return &Config{
        PostgresURL: postgresURL,
        JWTKey:     jwtKey,
    }, nil
}