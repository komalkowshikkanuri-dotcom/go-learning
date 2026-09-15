package main

import (
	"net/http"
	"fmt"
	"strings"
	"strconv"
)

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
	
		//fmt.Fprintln(w, "This is the user handle")
		parts := strings.Split(r.URL.Path, "/")

		userID := parts[len(parts)-1]

		userIDInt, err := strconv.Atoi(userID)

		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		} else {
			fmt.Fprintln(w, "User ID: ", userIDInt)
		}
	} else {
		http.Error(w, "Method not Allowed", 405)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from Go")
	fmt.Fprintln(w, "Method: ", r.Method)
	fmt.Fprintln(w, r.URL)
	fmt.Fprintln(w, r.URL.Path)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/users/", usersHandler)
	http.ListenAndServe(":8080", nil)
}