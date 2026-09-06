package main

import "fmt"

func main(){
	numbers := [5]int{12, 45, 7, 89, 23}

	maxNum := 0
	maxIndex := 0

	for index, num := range numbers{
		if(num > maxNum){
			maxNum = num
			maxIndex = index
		}
	}

	fmt.Println(maxIndex, maxNum)
}