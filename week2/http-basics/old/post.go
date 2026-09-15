package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users []User
var nextID = 1
var mu sync.Mutex

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		mu.Lock()

		data, err := json.Marshal(users)
	
		mu.Unlock()

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "application/json")
		fmt.Fprintln(w, string(data))
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if user.Age < 1 || user.Age > 100 || user.Name == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	mu.Lock()

	user.ID = nextID
	nextID++
	users = append(users, user)

	mu.Unlock()

	fmt.Println(user)

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(user)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, string(data))
}

func main() {
	http.HandleFunc("/users", createUserHandler)
	http.ListenAndServe(":8080", nil)
}