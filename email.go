package beasty

// import (
// 	"net/smtp"
// )

type email struct {
	content        string
	destination    string
	sender         string
	defaultHeaders struct{}
	optHeaders     struct{}
}

func newEmail(d string) email {
	e := email{}
	return e
}

func (e email) setContent(c string) {
	e.content = c
}

func (e email) setDestination(d string) {
	e.destination = d
}

func (e email) setSender(s string) {
	e.sender = s
}

func (e email) setOptHeaders(headers struct{}) {
	e.optHeaders = headers
}
