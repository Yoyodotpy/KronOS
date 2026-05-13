package main

type Winfo struct {
	ID     int
	Width  int
	Height int
	Stride int
	Pixels []byte
}

type WindowReply struct {
	Success bool
	Focused bool
}

type Title struct {
	Title string
	ID    int
}
