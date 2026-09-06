package main

import "fmt"

func main(){
	capital := map[string]string{
			"India": "New Delhi",
			"France": "Paris",
			"Japan": "Tokyo",
			"Germany": "Berlin",
		}

	_, USAExists := capital["USA"]
	_, IndiaExists := capital["India"]
	
	fmt.Println("India exists: ", IndiaExists)
	fmt.Println("USA exists:", USAExists)

	for country, city := range capital {
		fmt.Println(country, "-->", city)
	}
}