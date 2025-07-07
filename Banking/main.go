package main

import "log"

func main() {
	db := initDB()
	defer db.Close()
	// Starting port at 1000
	
	store := &PostgresStore{db: db}
	server := NewAPIServer(":1000", store)
	log.Println("Server starting on port 1000...")
	if err := server.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	
	}
}
