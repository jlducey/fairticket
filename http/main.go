package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type logWriter struct{} // our struct is empty but allows us to make functions

func main() {
	resp, err := http.Get("http://www.google.com")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	//bs := make([]byte, 99999) // make a big byte slice to read data in to
	//resp.Body.Read(bs)        // take resp returned above and read the body in to the byte slice
	//fmt.Println(string(bs))   // print byte slice out, after converting to string.. send to more to see html

	//io.Copy(os.Stdout, resp.Body) // copies resp.Body to writer interface impleting code aimed to stdout

	lw := logWriter{}      // we create a logwriter
	io.Copy(lw, resp.Body) // we use our logwriter to handle resp.Body with its function
}

// note return type in following function, int and error as required by the reader/writers we are connecting to
func (logWriter) Write(bs []byte) (int, error) { //reader/Writer interface required we process resp.Body and return int of bytes and err
	fmt.Println(string(bs))                             // our byte string from resp.Body is converted to string and printed
	fmt.Println("Just wrote this many bytes:", len(bs)) // we can see how many bytes by getting length of the byte string
	return len(bs), nil                                 // we return the # of bytes read and the err as required by the reader/writer interfaces
}
