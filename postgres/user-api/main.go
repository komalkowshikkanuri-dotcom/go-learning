package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/lib/pq"
	//_ "github.com/lib/pq"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, age FROM users")

		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var users []User

		for rows.Next() {
			var user User

			err = rows.Scan(
				&user.ID,
				&user.Name,
				&user.Age,
			)

			if err != nil {
				http.Error(w, "Database Error", http.StatusInternalServerError)
				return
			}

			users = append(users, user)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(users)

		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func getUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		userIDInt, err := strconv.Atoi(parts[len(parts)-1])

		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		var user User

		err = db.QueryRow("SELECT id, name, age FROM users WHERE id = $1", userIDInt).Scan(&user.ID, &user.Name, &user.Age)

		if err == sql.ErrNoRows {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func createUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var user User
		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if user.ID <= 0 ||
			user.Name == "" ||
			user.Age <= 0 ||
			user.Age > 100 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		_, err = db.Exec(
			"INSERT INTO users (id, name, age) VALUES ($1, $2, $3)",
			user.ID, user.Name, user.Age,
		)

		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" {
					http.Error(w, "User already exists", http.StatusConflict)
					return
				}
			}

			log.Println(err)
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(user)

		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, string(data))
		return
	}
}

func main() {
	db, err := sql.Open(
		"postgres",
		"host=localhost port=5432 user=postgres dbname=backend_db sslmode=disable",
	)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL")

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUsers(db)(w, r)
		case http.MethodPost:
			createUser(db)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/users/", getUser(db))

	log.Println("Server running on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
