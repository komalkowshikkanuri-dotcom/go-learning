package main

import "fmt"

func main(){
	words := []string{
    		"apple",
    		"banana",
    		"apple",
    		"orange",
    		"banana",
    		"apple",
	}

	count := make(map[string]int)
	maxLengthWord := ""

	for _, word := range words {
		count[word]++
		if(len(word) > len(maxLengthWord)) {
			maxLengthWord = word
		}
	}

	for word, val := range count{
		fmt.Println(word, "->", val)
	}

	fmt.Println("Biggest word: ",maxLengthWord)
}