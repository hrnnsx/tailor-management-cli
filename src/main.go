package main

import (
	"fmt"
	"log"
	"tailor-management-cli/cli"
	"tailor-management-cli/config"

	// "tailor-management-cli/entity"

	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
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

	prompt := promptui.Select{
		Label: "Pilih menu",
		Items: cli.SignMenu,
	}
	index, _, err := prompt.Run()

	switch index {
	case 0:
		// signIn()
	case 1:
		// signUp()
	case 2:
		return
	}
}
