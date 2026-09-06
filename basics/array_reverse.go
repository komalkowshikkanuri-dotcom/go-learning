package main

import "fmt"


func main(){
	numbers := [5]int{10, 20, 30, 40 , 50}

	for i := range len(numbers) / 2 {
		j := len(numbers) - 1 - i
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	
	for _, num := range numbers{
		fmt.Println(num)
	}
}