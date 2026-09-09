package main

import (
	"fmt"
	"os"
)

type User struct {
	Age  int
	Name string
}

func saveUser(user User) error {
	data := fmt.Sprintf("Name: %s, \nAge: %d", user.Name, user.Age)
	err := os.WriteFile("user.txt", []byte(data), 0644)
	
	if err != nil{
		return err
	}
	
	return nil
}

func main() {
	user := User{
		Name: "Alex",
		Age: 26,
	}

	err := saveUser(user)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("User Saved Successfully")

	
	data, err := os.ReadFile("user.txt")
	
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	
	str := string(data)
	fmt.Println(str)
}