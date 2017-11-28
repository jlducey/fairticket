package main

import (
	"fmt"
	"net/http"
)

func main() {
	links := []string{
		"http://www.amazon.com",
		"http://www.facebook.com",
		"http://golang.org",
		"http://stackoverflow.com",
		"http://www.google.com",
	}
	c := make(chan string) // make a channel c which only accepts string messages

	for _, link := range links {
		go checkLink(link, c) // pass the channel c to functions like any other variable
	}
	for func(){
		time.Sleep(5 *time.Second)
	    checkLink(l, c) // listen for messages on the channel
	}()
}

func checkLink(link string, c chan string) { // pass in the channel giving name chan type
	_, error := http.Get(link)
	if error != nil {
		fmt.Println(link, "might be down")
		//c <- "might be down" // send message on channel if down
		c <- link
		return
	}
	fmt.Println(link, "looks to be up")
	//c <- "guess its up" // send message on channel if up
	c <- link
}
