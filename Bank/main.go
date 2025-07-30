package main

import (
	"Bank/controllers"
	"Bank/db"
	"log"
)

func main() {
	database := db.InitDB()
	defer database.Close()

	storage := &db.PostgresStore{
		DB: database,
	}

	// Create tables
	if err := storage.CreateAccountTable(); err != nil {
		log.Fatalf("Failed to create account table: %v", err)
	}
	if err := storage.CreateUsersTable(); err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Start the API server
	server := controllers.NewAPIServer(":1000", storage)
	log.Println("Server starting on port 1000...")
	if err := server.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
