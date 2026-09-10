package main

import (
	"fmt"
	"sync"
)

func sendMessage(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()	
	
	value := <-ch
	fmt.Println(value)
}

func main() {
	ch := make(chan int)
	
	var wg sync.WaitGroup
	wg.Add(1)

	go sendMessage(ch, &wg)

	ch <- 42

	wg.Wait()
}