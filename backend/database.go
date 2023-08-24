package main

import (
	"database/sql"
	"log"
	"os"
)

var db *sql.DB

func initializeDatabase() {
	var err error

	// Get database details from environment variables
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	// Create connection string from environment variables
	connectionString := "host=" + host + " dbname=" + dbname + " user=" + user + " password=" + password + " sslmode=disable"

	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS people (
            id SERIAL PRIMARY KEY,
            name TEXT,
            age INT,
            address TEXT
        )
    `)
	if err != nil {
		log.Fatal(err)
	}
}
