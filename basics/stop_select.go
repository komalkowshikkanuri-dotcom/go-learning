package main

import (
	"fmt"
)

func sendNumber(ch chan<- int) {
	ch <- 10
	ch <- 20
	ch <- 30

	ch <- 40

	close(ch)
}

func main() {
	ch := make(chan int)

	go sendNumber(ch)
	
	for{
		select {
		case value, ok := <-ch:
			if !ok {
				return
			} 					

			fmt.Println(value)
				
		}
	}	
}