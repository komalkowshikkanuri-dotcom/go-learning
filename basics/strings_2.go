package main

import "fmt"

func main(){
	var word string
	
	fmt.Println("Enter a word: ")
	fmt.Scan(&word)

	var count int
	
	for _, char := range word {
		if char == 'a' || char == 'e' || char == 'i' || char == 'o' || char == 'u' {
			count++
		}
	}

	fmt.Println(count)
}