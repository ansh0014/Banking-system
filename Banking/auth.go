package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a new JWT token using the configured secret key
func GenerateJWTWithUsername(username string) (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load configuration: %v", err)
	}
	secretKey := []byte(config.JWTKey)
	if len(secretKey) == 0 {
		return "", fmt.Errorf("JWT key is not configured")
	}

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateJWT parses and validates a JWT token using the configured secret key
func ValidateJWT(tokenString string) error {
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}
	secretKey := []byte(config.JWTKey)
	if len(secretKey) == 0 {
		return fmt.Errorf("JWT key is not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is as expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Token is valid!")
		fmt.Printf("Claims: %v\n", claims)
		return nil
	}

	return fmt.Errorf("invalid token")
}
