package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {

	dbString := os.Getenv("DB_URL")
	dbURL, err := sql.Open("mysql", dbString)
	if err != nil {
		return nil, fmt.Errorf("Gagal membuka koneksi database: %w", err)
	}

	if err := dbURL.Ping(); err != nil {
		return nil, fmt.Errorf("Gagal ping database: %w", err)
	}

	fmt.Println("Successfully connected to database :D")
	return dbURL, nil
}
