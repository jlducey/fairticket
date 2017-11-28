package main

func main() {
	//cards := newDeckFromFile("my_cards")
	cards := newDeck()
	//hand, remainingDeck := deal(cards, 5)
	cards.shuffle()
	cards.print()
	//hand.print()
	//remainingDeck.print()
	//fmt.Println(cards.toString())
	//cards.saveToFile("my_cards") // saves to disk our deck
}
