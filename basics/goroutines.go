package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello from goroutine!")
}

func printNumbers() {
	for i := 1; i <= 5; i++ {
		fmt.Println("Number: ", i)
		time.Sleep(500 * time.Millisecond)
	}
}

func printLetters() {
	for _, letter := range []string{"A", "B", "C", "D", "E"} {
		fmt.Println("Letters: ", letter)
		time.Sleep(500 * time.Millisecond)
	}
	
}

func main() {
	go sayHello()
	time.Sleep(200 * time.Millisecond)
	
	fmt.Println("Main finished")

	go printNumbers()
	go printLetters()

	time.Sleep(5000 * time.Millisecond)
}