package main

type storage struct {
	redis struct{}
}

func NewStorage() *storage {
	s := &storage{}
	return s
}
