package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strings"
	"strconv"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func userHandler(w http.ResponseWriter, r *http.Request, user User) {
	// Check if the request is GET or not
	if r.Method != "GET" {
		http.Error(w, "Invalid Input", http.StatusMethodNotAllowed)
		return
	}

	//Converting the string into parts and then using the last part as the userID and 
	//Printing the error message if there is one
	parts := strings.Split(r.URL.Path, "/")
	userID, err := strconv.Atoi(parts[len(parts) - 1])

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	
	//Comparing the User ID in URL with the user of that of User data type
	if user.ID != userID {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	//format of response body
	w.Header().Set("Content-Type", "application/json")
	
	//converting to Json 
	//Printing an error message if there is one
	data, err := json.Marshal(user)

	if err != nil {
		http.Erro(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//Printing JSON data
	fmt.Fprintln(w, string(data))
}

func main() {
	user := User {
		ID: 67,
		Name: "David",
		Age: 32,
	}

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
			userHandler(w, r, user)
		})

	http.ListenAndServe(":8080", nil)
}