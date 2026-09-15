package main

import (
	"fmt"
	"net/http"
	"strings"
	"strconv"
	"encoding/json"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Welcome to the page")

	fmt.Fprintln(w, r.Method)
	fmt.Fprintln(w, r.URL.Path)
}

func userHandler(w http.ResponseWriter, r *http.Request, user User) {
	if r.Method == "GET" {
		parts := strings.Split(r.URL.Path, "/")

		userIDInt, err := strconv.Atoi(parts[len(parts) - 1])

		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return 
		}

		if user.ID != userIDInt {
    			http.Error(w, "User not found", http.StatusNotFound)
    			return
		}

		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(user)
		
		if err != nil {
			fmt.Fprintln(w, "Trouble converting data to JSON")
		} else {			
			fmt.Fprintln(w, string(data))
		}
	} else {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}
}

func main() {
	user := User {
			ID: 25,
			Name: "Alice",
			Age: 20,
		}

	http.HandleFunc("/", handler)
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		userHandler(w, r, user)
	})

	http.ListenAndServe(":8080", nil)
}