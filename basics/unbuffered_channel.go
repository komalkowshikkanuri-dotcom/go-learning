package main

import "fmt"

func sendNumbers(ch chan<- int) {
	ch <- 10
	ch <- 20
	ch <- 30

	close(ch)
}

func main() {
	ch := make(chan int)
	
	go sendNumbers(ch)
	
	for value := range ch {
		fmt.Println(value)
		
	}
} 