package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, error := os.Open(os.Args[1])
	if error != nil {
		fmt.Println("error was:", error)
	}
	io.Copy(os.Stdout, file)
}
