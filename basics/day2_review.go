package main

import "fmt"

func main(){
	arr := [5]int{5, 10, 15, 20, 25}
	fmt.Println(arr[2])


	slice := []int{10, 20, 30}
	slice = append(slice, 40, 50)
	fmt.Println(slice)


	capitals := map[string]string{
		"India": "New Delhi",
		"Japan": "Tokyo",
		"France": "Paris",
	}
	fmt.Println(capitals["France"])
	

	_, ok := capitals["Germany"]
	if ok{
		fmt.Println("Germany Exists!!!")
	} else{
		fmt.Println("Germany doesn't exist")
	}

	
	numbers := []int{10, 20, 30, 40, 50}
	for _, val := range numbers {
		fmt.Println(val)
	}

	
	word := "hello"
	for _, char := range word {
		fmt.Println(string(char))
	}

	
	ans := 0		
	number := []int{4, 8, 15, 16, 23, 42}
	for _, val := range number{
		ans += val
	}
	fmt.Println(ans)
}