package main

import (
	"net/http"
	"fmt"
)

func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is the user handle")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from Go")
	fmt.Fprintln(w, "Method: ", r.Method)
	fmt.Fprintln(w, r.URL)
	fmt.Fprintln(w, r.URL.Path)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/users", usersHandler)
	http.ListenAndServe(":8080", nil)
}