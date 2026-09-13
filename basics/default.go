package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	
	go func() {
		time.Sleep(2 * time.Millisecond)
		ch <- 100
	}()

	for {
		select {
			case value := <-ch:
				fmt.Println("The value is:", value)
				return

			default:
				fmt.Println("Waiting...")

		}
	}
}