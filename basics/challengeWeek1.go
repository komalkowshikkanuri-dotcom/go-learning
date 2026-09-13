package main

import (
	"fmt"
	"sync"
)

type Result struct {
	ID     int
	Square int
}

func workers(numbers <-chan int, id int, wg *sync.WaitGroup, results chan<- Result) {
	defer wg.Done()

	for number := range numbers {
		results <- Result{
			ID: id,
			Square : number * number,
		}
	}
}

func main() {
	var wg sync.WaitGroup

	numbers := make(chan int)
	results := make(chan Result)

	wg.Add(3)

	go workers(numbers, 1, &wg, results)
	go workers(numbers, 2, &wg, results)
	go workers(numbers, 3, &wg, results)

	go func() {
		numbers <- 10
		numbers <- 5
		numbers <- 3
		numbers <- 6
		numbers <- 8

		close(numbers)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for value := range results{
		fmt.Println("Worker: ", value.ID, ":", value.Square)
	}
}