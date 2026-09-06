package main

import "fmt"

func isEven(a int) bool{
	if a % 2 == 0{
		return true
	} else{
		return false
	}
}

func multiply(a, b int) int{
	return a * b
}

func main(){
	ans := multiply(6, 7)
	
	fmt.Println(ans)

	even := isEven(10)

	if even{
		fmt.Println("The number is even")
	}else{
		fmt.Println("The number is odd")
	}

	que, rem := divide(17, 5)

	fmt.Println("Quotient: ", que)
	fmt.Println("Reminder: ", rem)
}

func divide(a, b int) (int, int) {
	return a / b, a % b
}