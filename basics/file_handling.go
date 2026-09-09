package main

import (
	"fmt"
	"os"
)

func main() {
	err := os.WriteFile("hello.txt", []byte("Hello from Go"), 0644)

	if err != nil{ 
		fmt.Println("Error Writing file", err)
		return
	}

	fmt.Println("File written successfully")

	data, err := os.ReadFile("hello.txt")

	if err != nil {
		fmt.Println("Error Reading File", err)
		return
	}
	
	str := string(data)
	fmt.Println(str)
}