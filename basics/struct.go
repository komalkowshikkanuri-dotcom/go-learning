package main

import "fmt"

type User struct{
		Name string
		Age int
		Email string
	}

func main(){
	user := User{
		Name: "Alex",
		Age: 25,
		Email: "alex@example.com",
	}

	user1 := User{
		Name: "Bob",
		Age: 26,
		Email: "bob@example.com",
	}

	
	fmt.Println("Name: ", user.Name)
	fmt.Println("Age: ", user.Age)
	fmt.Println("Email: ", user.Email)

	fmt.Println("Age: ", user1.Age)

	user.birthday()
	user1.birthday()
	
	message := user.introduce()
	fmt.Println(message)
	message1 := user1.introduce()
	fmt.Println(message1)
}

func (user User) introduce() string{
	return fmt.Sprintf("Hi, I'm %s and I'm %d years old", user.Name, user.Age)
}

func (user *User) birthday(){
	user.Age = user.Age + 1
}