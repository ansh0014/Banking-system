package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var input User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Username == "" || input.Password == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var stored User
	// TODO: Implement proper database query once DB is set up
	// For now, this is a placeholder
	if input.Username != "testuser" {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	stored.Username = input.Username
	stored.Password = "$2a$10$examplehashedpassword1234567890abcdef" // placeholder hashed password

	if err := bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Use GenerateJWT function with username
	token, err := GenerateJWT(stored.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// GenerateJWT creates a JWT token for a specific username
func GenerateJWT(username string) (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}
	secretKey := []byte(config.JWTKey)
	if len(secretKey) == 0 {
		return "", fmt.Errorf("JWT key is not configured")
	}

	expiration := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
