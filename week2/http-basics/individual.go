package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users []User

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.PathValue("id")

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	found := false
	var foundIndex int

	for index, value := range users {
		if value.ID == userIDInt {
			found = true
			foundIndex = index
			break
		}
	}

	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	users = append(users[:foundIndex], users[foundIndex+1:]...)
	fmt.Fprintln(w, "Successfully Removed the user")

	fmt.Fprintln(w, users)

	return
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.PathValue("id")

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	found := false
	var foundUser User

	for _, value := range users {
		if value.ID == userIDInt {
			found = true
			foundUser = value
			break
		}
	}

	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	data, err := json.Marshal(foundUser)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(data))
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.PathValue("id")

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var updatedUser User

	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if updatedUser.Age < 1 || updatedUser.Age > 100 || updatedUser.Name == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	found := false

	for index, value := range users {
		if value.ID == userIDInt {
			users[index].Name = updatedUser.Name
			users[index].Age = updatedUser.Age
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	data, err := json.Marshal(updatedUser)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(data))
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth != "secret123" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth != "secret123" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users/{id}", getUserHandler)
	mux.HandleFunc("PUT /users/{id}", updateUserHandler)
	mux.HandleFunc("DELETE /users/{id}", deleteUserHandler)
	
	return mux
}

func main() {
	users = []User{
		{ID: 1, Name: "John", Age: 25},
		{ID: 2, Name: "Alice", Age: 30},
	}

	mux := setupRoutes()
	handler := loggingMiddleware(authMiddleware(mux))

	http.ListenAndServe(":8080", handler)
}