package main

import (
	"fmt"
	"log"
	"tailor-management-cli/config"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := config.Connect()
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer db.Close()

	fmt.Println("Cuayo")
}
