package main

import (
	"Bank/Authantication" // Fixed spelling
	"Bank/controllers"
	"Bank/db"
	"fmt"
	"log"
)

func main() {
	database := db.InitDB()
	defer database.Close()

	storage := &db.PostgresStore{
		DB: database, // Fix: field name should be 'Db' (exported field)
	}

	// Generate JWT token
	token, err := auth.GenerateJWT("exampleUser") // Fix: use correct package name
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return
	}
	fmt.Printf("Generated Token: %s\n", token)

	// Validate the JWT
	err = auth.ValidateJWT(token) // Fix: use correct package name
	if err != nil {
		fmt.Printf("Error validating token: %v\n", err)
	} else {
		fmt.Println("Token validation successful!")
	}

	// Create tables
	if err := storage.CreateAccountTable(); err != nil {
		log.Fatalf("Failed to create account table: %v", err)
	}
	if err := storage.CreateUsersTable(); err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	server := controllers.NewAPIServer(":1000", storage)
	log.Println("Server starting on port 1000...")
	if err := server.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
