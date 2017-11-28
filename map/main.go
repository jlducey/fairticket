package main

import "fmt"

func main() {
	colors := map[string]string{ // can make map this way or way of following line
		//colors := make(map[int]string)
		"red":   "27abj", // key value pairs can be done this way, even last pair gets comma after
		"green": "28xxe",
		"blue":  "oxy33",
	}
	//colors[10] = "27abj" //alternative perl like way to assign values.. note map is typed, so only int and strings here
	//fmt.Println(colors)
	//delete(colors, 10) // how you remove a value from map.. its like a perl hash pretty much so far
	//fmt.Println(colors)
	printMap(colors)
}
func printMap(c map[string]string) {
	for key, value := range c { // key and value at same time, cool
		fmt.Printf("key is:%v and value is :%v\n", key, value)
	}
}
