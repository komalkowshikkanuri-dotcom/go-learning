package main

import (
	"fmt"
	"errors"
)

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	
	return a / b, nil
}

func handleResult(result int, err error) {
	if err != nil {
		fmt.Println("Some error has occurred: ", err)
	} else{
		fmt.Println(result)
	}
}

func main() {
	result, err := divide(10, 0)

	handleResult(result, err)
	
	result, err = divide(10, 2)

	handleResult(result, err)
}