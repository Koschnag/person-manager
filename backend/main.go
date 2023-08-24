package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Model
type Person struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}

var db *sql.DB

// Database Initialization
func initializeDatabase() {
	var err error
	connectionString := "host=localhost dbname=person_management user=yourUser password=yourPassword sslmode=disable"
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

// Handlers
func GetPeople(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, age, address FROM people")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var people []Person
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Age, &p.Address); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		people = append(people, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

// Main function
func main() {
	initializeDatabase()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/people", GetPeople).Methods("GET")
	// You can add more routes for POST, PUT, DELETE...

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
