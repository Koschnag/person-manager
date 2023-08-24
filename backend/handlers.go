package main

import (
	"encoding/json"
	"net/http"
)

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

// Additional CRUD handlers for adding, updating, and deleting people can be added here.
