package main

import "fmt"

type shape interface { // interface wants a getArea function defined with a float64 retured with area
	getArea() float64
}
type square struct {
	sideLength float64
}
type triangle struct {
	height float64
	base   float64
}

func main() {
	s := square{2.0}
	t := triangle{2.0, 3.0}
	printArea(s)
	printArea(t)
}

func (s square) getArea() float64 {
	return s.sideLength * s.sideLength
}

func (t triangle) getArea() float64 {
	return 0.5 * t.base * t.height
}

func printArea(s shape) {
	fmt.Println(s.getArea())
}
