package main

import (
	"fmt"
	"time"
	"sync"
)

func main() {
	fastResults := make(chan int, 1)
	slowResults := make(chan int, 1)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		fastResults <- 100
		
	}()

	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		slowResults <- 200
	}()
	
	for i := 0; i < 2; i++ {
		select {
			case value := <-fastResults:
				fmt.Println("Fast Result", value)

			case value := <-slowResults:
				fmt.Println("Slow Result", value)
		}
	}

	wg.Wait()
}