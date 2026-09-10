package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func ConnectDB() *sql.DB {

	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Berkas .env tidak ditemukan, menggunakan environment sistem.")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	fmt.Println("Koneksi database berhasil.")
	return db
}
