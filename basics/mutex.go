package main

import (
	"fmt"
	"sync"
)

func increment(counter *int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	defer mu.Unlock()

	for i := 1; i <= 1000; i++ {
		*counter += 1
	}
}

func main() {
	counter := 0

	var wg sync.WaitGroup
	var mu sync.Mutex
	
	wg.Add(2)
	go increment(&counter, &mu, &wg)
	go increment(&counter, &mu, &wg)

	wg.Wait()

	fmt.Println(counter)
}