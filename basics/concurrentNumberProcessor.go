package main

import (
	"fmt"
	"sync"
)

type Print struct {
	ID     int
	Square int
}

func worker(id int, numbers <-chan int, wg *sync.WaitGroup, result chan<- Print) {
	defer wg.Done()

	for number := range numbers {
    		result <- Print{
				ID: id,
				Square: number * number,
			}
		
	}
}

func main() {
	numbers := make(chan int)
	result  := make(chan Print)

	var wg sync.WaitGroup

	wg.Add(3)

	go worker(1, numbers, &wg, result)
	go worker(2, numbers, &wg, result)
	go worker(3, numbers, &wg, result)
	
	go func() {
		numbers <- 10
		numbers <- 5
		numbers <- 3
		numbers <- 6
		
		close(numbers)
	} ()

	go func() {
		wg.Wait()
		close(result)
	} ()

	for value := range result {
		fmt.Println("Worker", value.ID, ":", value.Square)
	}	
}