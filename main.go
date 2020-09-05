package main

import (
	"os"
)

// main driver for program
// initializes bot and sets to listen for guild events
func main() {
	var token string
	// if cli arguments passed, assume ENV override
	if len(os.Args) != 1 {
		arguments := os.Args[1:]
		token = arguments[0]
	} else {
		token = os.Getenv("TOKEN")
	}
	b := NewBeasty(token)
	b.Start()
}
