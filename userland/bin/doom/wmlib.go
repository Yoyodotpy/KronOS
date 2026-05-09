package main

type Winfo struct {
	ID     int
	Width  int
	Height int
	Stride int
	Pixels []byte
	Title  string
}

type WindowReply struct {
	Success bool
}
