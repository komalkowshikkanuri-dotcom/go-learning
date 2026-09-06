package main

import "fmt"

func main(){
	numbers := []int{10, 20, 30}
	fmt.Println("Before:", len(numbers), cap(numbers))

	numbers = append(numbers, 40, 50, 60, 70)
	
	fmt.Println(len(numbers))
	fmt.Println(cap(numbers))
}