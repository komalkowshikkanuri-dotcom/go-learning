package main

import (
	"fmt"	
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		ch1 <- "Hello from 1"
	}()
	
	go func() {
		ch2 <- "Hello from 2"
	}()

	time.Sleep(100 * time.Millisecond)
	
	select {
		case msg := <-ch1:
    			fmt.Println(msg)
		case msg := <-ch2:
    			fmt.Println(msg)
		default:
    			fmt.Println("No messages available")
	}
}      