package main

import (
	"tailor-management-cli/cli"
	"tailor-management-cli/cli/auth"
	"tailor-management-cli/config"
)

func main() {
	// Koneksi database (godotenv dan pengecekan error sudah ditangani di dalam ConnectDB)
	db := config.ConnectDB()
	defer db.Close()

	// Memanggil menu utama CLI
	cli.MainMenu(db)
}
