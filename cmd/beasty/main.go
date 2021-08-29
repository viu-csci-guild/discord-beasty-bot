package main

import (
	"os"

	"github.com/viu-csci-guild/beasty/cmd/beasty/client"
)

func main() {

	// if cli arguments passed, assume ENV override
	token := os.Getenv("TOKEN")
	studentRoleID := os.Getenv("STUDENT_ROLE_ID")
	startupRoomID := os.Getenv("START_ROOM_ID")
	serverID := os.Getenv("SERVER_ID")
	squelch_check := os.Getenv("SQUELCH_STARTUP")
	squelch := false
	if squelch_check == "yes" {
		squelch = true
	}

	b := client.NewBeasty(token, studentRoleID, startupRoomID, serverID, squelch)
	b.Start()
}
