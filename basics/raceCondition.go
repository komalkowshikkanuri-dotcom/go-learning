package main

import (
	"fmt"
	"sync"
)

func increment(counter *int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= 1000; i++ {
		*counter += 1
	}
}

func main() {
	counter := 0

	var wg sync.WaitGroup
	
	wg.Add(2)
	go increment(&counter, &wg)
	go increment(&counter, &wg)

	wg.Wait()

	fmt.Println(counter)
}