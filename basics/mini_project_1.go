package main

import "fmt"

type User struct {
	Name  string
	Age   int
	Email string
}

func (user User) introduce() string {
	return fmt.Sprintf("Hi, I'm %s and I'm %d years old", user.Name, user.Age)
}

func (user *User) birthday() {
	user.Age += 1
}

type Person interface {
	introduce() string
}

func greet(person Person) string{
	return person.introduce()
}

func main() {
	user1 := User {
			Name: "Alex", 
			Age: 25, 	
			Email: "Alex@example.com",
		}

	user2 := User {
			Name: "Bob",
			Age: 30,
			Email: "bob@example.com",
		}

	fmt.Println(greet(user1))
	fmt.Println(greet(user2))

	user1.birthday()
	fmt.Println(greet(user1))
	
}