package main

import (
	"fmt"
	"sync"
	"time"
)

func sayHello() {
	fmt.Println("Hello from goroutine!")
}

func printNumbers(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= 5; i++ {
		fmt.Println("Number: ", i)
		time.Sleep(500 * time.Millisecond)
	}
}

func printLetters(wg *sync.WaitGroup) {
	defer wg.Done()

	for _, letter := range []string{"A", "B", "C", "D", "E"} {
		fmt.Println("Letters: ", letter)
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	go sayHello()
	time.Sleep(200 * time.Millisecond)

	fmt.Println("Main finished")

	var wg sync.WaitGroup

	wg.Add(2)

	go printNumbers(&wg)
	go printLetters(&wg)

	wg.Wait()
}
