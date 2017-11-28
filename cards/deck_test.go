package main

import (
	"os"
	"testing"
)

func TestNewDeck(t *testing.T) { // handler by convention, gives type
	d := newDeck()    // we make a new deck d
	if len(d) != 52 { // should have 52 cards in a deck else error
		t.Errorf("Expected 52 but got: %v", len(d)) // note E is capitalized and it needs printf style format string
	}
	if d[0] != "Ace of Spades" { // check first card in deck if not right one error
		t.Errorf("Expected Ace of Spades but got: %v", d[0])
	}
	if d[51] != "King of Clubs" { //check last card note index is 51 due to count starting at 0
		t.Errorf("Expected King of Clubs but got: %v", d[51])
	}
}

func TestSaveToDeckAndNewDeckFromFile(t *testing.T) {
	os.Remove("_decktesting")
	deck := newDeck()
	deck.saveToFile("_decktesting")
	loadedDeck := newDeckFromFile("_decktesting")
	if len(loadedDeck) != 52 {
		t.Errorf("Expected 52 cards in deck, but got: %v", len(loadedDeck))
	}
	os.Remove("_decktesting")
}
