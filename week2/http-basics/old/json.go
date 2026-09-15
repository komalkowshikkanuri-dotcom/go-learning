package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	user := User {
		ID: 25,
		Name: "Komal",
		Age: 21,
	}

	data, err := json.Marshal(user)

	if err != nil {
		fmt.Println("Some error has occured while converting the data into JSON")
	} else {
		fmt.Println(string(data))
	}
}