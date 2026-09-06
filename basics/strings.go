package main

import "fmt"

func main(){
	word := "hello"

	fmt.Println("First letter:", string(word[0]))
	fmt.Println("Last letter:", string(word[len(word) - 1]))

	fmt.Println("Word length:", len(word))

	for _, char := range word {
		fmt.Println(string(char))
	}
}