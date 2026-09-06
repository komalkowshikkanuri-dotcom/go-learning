package main

import "fmt"

func main(){
	numbers := []int{10, 20, 30, 40, 50}

	numbers = append(numbers, 60)
	numbers = append(numbers, 70)
	numbers = append(numbers, 80)


	fmt.Println(numbers[1:4])
	fmt.Println(numbers[5:])
}