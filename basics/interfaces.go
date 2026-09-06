package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
}

type Circle struct {
	radius float64
}

type Rectangle struct {
	width  float64
	height float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func totalArea(shape Shape) float64 {
	return shape.Area()
}


func main() {
	c := Circle{radius: 5.0}
	r := Rectangle{width: 5.0, height: 4.0}

	fmt.Println(totalArea(c))
	fmt.Println(totalArea(r))
}


