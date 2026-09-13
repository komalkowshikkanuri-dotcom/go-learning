package main

import "fmt"

func goRoutine1(ch1 chan<- string) {
	ch1 <- "Message 1"
	ch1 <- "Message 2"

	close(ch1)
}

func goRoutine2(ch2 chan<- string) {
	ch2 <- "Message A"
	ch2 <- "Message B"

	close(ch2)
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go goRoutine1(ch1)
	go goRoutine2(ch2)

	for ch1 != nil || ch2 != nil {
		select {
			case msg, ok := <- ch1:
				if ok {
					fmt.Println(msg)
				} else {
					ch1 = nil
				}
	
			case msg, ok := <- ch2:
				if ok {
					fmt.Println(msg)
				} else {
					ch2 = nil
				}
			
		}

	}
}