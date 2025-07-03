package main

import "log"

func main() {
	db := initDB()
	defer db.Close()
	
	server := NewAPIServer(":1000")
	log.Println("Server starting on port 1000...")
	if err := server.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	
	}
}
