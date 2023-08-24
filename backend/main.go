package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Person struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}

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

func initializeDatabase() {
	// Get database details from environment variables
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	// Create connection string from environment variables
	connectionString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", host, dbname, user, password)

	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS people (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            age INT NOT NULL,
            address TEXT
        )
    `)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

// Create a new person
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var p Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO people (name, age, address) VALUES ($1, $2, $3)", p.Name, p.Age, p.Address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Get a single person by ID
func GetPerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	row := db.QueryRow("SELECT id, name, age, address FROM people WHERE id = $1", id)

	var p Person
	if err := row.Scan(&p.ID, &p.Name, &p.Age, &p.Address); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// Update an existing person
func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var p Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE people SET name=$1, age=$2, address=$3 WHERE id=$4", p.Name, p.Age, p.Address, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete a person
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec("DELETE FROM people WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	initializeDatabase()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/people", GetPeople).Methods("GET")
	r.HandleFunc("/people", CreatePerson).Methods("POST")
	r.HandleFunc("/people/{id:[0-9]+}", GetPerson).Methods("GET")
	r.HandleFunc("/people/{id:[0-9]+}", UpdatePerson).Methods("PUT")
	r.HandleFunc("/people/{id:[0-9]+}", DeletePerson).Methods("DELETE")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
