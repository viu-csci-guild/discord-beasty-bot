package main

import (
	"os"
	"testing"
)

// TestSmokeBeastyConn Test if beasty can connect to discord during construction
func TestSmokeBeastyConn(t *testing.T) {
	token := os.Getenv("DISCORD_API_KEY")
	b := NewBeasty(token)
	if b == nil {
		t.Errorf("Failed to construct beasty")
	}
}
