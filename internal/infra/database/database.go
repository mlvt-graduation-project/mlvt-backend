package infra

import (
	"database/sql"
	"log"
	"os"
)

var DB *sql.DB

func ConnectDatabase() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	log.Println("Connected to the database successfully")
}
