package main

import (
	"os"
)

// main driver for program
// initializes bot and sets to listen for guild events
func main() {
	var token string
	var localUse bool

	// if cli arguments passed, assume ENV override
	if len(os.Args) != 1 {
		localUse = true
		arguments := os.Args[1:]
		token = arguments[0]
	} else {
		token = os.Getenv("DISCORD_API_KEY")
	}
	b := NewBeasty(token)
	b.SetLocalUse(localUse)
	b.Start()
}
