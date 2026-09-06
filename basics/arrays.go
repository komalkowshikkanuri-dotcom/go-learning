package main

import "fmt"

func main(){
	numbers := [5]int{1,2,3,4,5}

	for _, num := range numbers{
		fmt.Println(num)
	}
}