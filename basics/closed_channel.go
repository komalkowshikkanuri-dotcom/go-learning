package main

import "fmt"

func sendNumber(ch chan int) {
	ch <- 100

	close(ch)
}

func main() {
	ch := make(chan int)

	go sendNumber(ch)

	value, ok := <-ch

	fmt.Println(value, ok)

	value, ok = <-ch

	fmt.Println(value, ok)

}