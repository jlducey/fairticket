package main

import ( // note import is multi-line, no commas and surrounded by PARENS
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Create a new type of 'deck'
// which is a slice of strings

type deck []string

// create representation of full deck of 52 playing cards
func newDeck() deck { // no OO but go has way to tie functions to a type, so we add functions that deck types can do
	cards := deck{}

	cardSuits := []string{"Spades", "Diamonds", "Hearts", "Clubs"}
	cardValues := []string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"}

	for _, suit := range cardSuits { // loop through 4 suits 13 values to make deck of 52 cards
		for _, value := range cardValues {
			cards = append(cards, value+" of "+suit)
		}
	}

	return cards
}

// deal func
func deal(d deck, handSize int) (deck, deck) { // example of 2 return values a dealthand and remaining cards both type deck
	return d[:handSize], d[handSize:]
}

// print out a deck func
func (d deck) print() { //   reciever for deck.. note reciever is before the function name
	for i, card := range d {
		fmt.Println(i, card)
	}
}

// convert deck to string
func (d deck) toString() string {
	return strings.Join([]string(d), ",")
}

//write byte-string to file
func (d deck) saveToFile(filename string) error {
	return ioutil.WriteFile(filename, []byte(d.toString()), 0666)
}

// reverse from file back to a deck
func newDeckFromFile(filename string) deck {
	bs, err := ioutil.ReadFile(filename) // read byte-string from file
	if err != nil {
		// option 1 .. log error and return call to newDeck()
		// option 2 .. log error and exit program
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	s := strings.Split(string(bs), ",") // break string in to slice of string
	return (deck(s))                    // type conversion from slice of string to deck, possible because deck was slice of string to start with
}

// shuffle deck .. give it new random order
func (d deck) shuffle() {
	source := rand.NewSource(time.Now().UnixNano()) // timestamp with UnixNano gives int64 as seed
	r := rand.New(source)                           // new random number is generated using seed from timestamp
	for i := range d {
		//newPosition := rand.Intn(len(d) - 1)        // psuedo random number, was always same order
		newPosition := r.Intn(len(d) - 1)           // random order replaces index number in next line
		d[i], d[newPosition] = d[newPosition], d[i] // swap position with random order

	}

}
