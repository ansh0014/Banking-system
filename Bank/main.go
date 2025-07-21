package main

import (
	"Bank/Authantication"
	"Bank/db"
	// "Bank/config"
	"Bank/controllers"
	"fmt"
	"log"
)

func main() {
	db := db.initDB()
	defer db.Close()
	// Starting port at 1000
		token, err := auth.GenerateJWT("exampleUser")
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return
	}
	fmt.Printf("Generated Token: %s\n", token)

	// Validate the JWT
	err = auth.ValidateJWT(token)
	if err != nil {
		fmt.Printf("Error validating token: %v\n", err)
	} else {
		fmt.Println("Token validation successful!")
	}
	
	Storage := &db.PostgresStore{db: db}
	if err := Storage.CreateAccountTable(); err != nil {
		log.Fatalf("Failed to create account table: %v", err)
	}
	server := controllers.NewAPIServer(":1000", Storage)
	log.Println("Server starting on port 1000...")
	if err := server.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
