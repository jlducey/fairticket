package main

import "fmt"

type bot interface { // ties only types together that have a getGreeting function and return strings
	getGreeting() string
}
type englishBot struct{}
type spanishBot struct{}

func main() {
	eb := englishBot{}
	sb := spanishBot{}
	printGreeting(eb) // calls interface passing englishbot
	printGreeting(sb) // calls interface passing spanishbot
}

func printGreeting(b bot) {
	fmt.Println(b.getGreeting()) // calls interface genericaly instead of specific get greeting
}

func (englishBot) getGreeting() string {
	// very custom fancy code
	return "Hi There"
}

func (spanishBot) getGreeting() string {
	//very spanish logic
	return "Hola"
}
