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

	counts := make(map[string]int)

	for _, word := range words {
		counts[word]++
	}

	for item, val := range counts{
		fmt.Println(item, "-->", val)
	}
}