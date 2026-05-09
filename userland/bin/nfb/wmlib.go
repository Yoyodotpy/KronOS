package main

import (
	"image"
	"os"
)

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

type window struct {
	ID     int
	X      int
	Y      int
	width  int
	height int
	img    image.RGBA
	title  string
}

func (win window) close() {
	proc, _ := os.FindProcess(win.ID)
	proc.Kill()
	delete(windows, win.ID)
	windowsneedrefresh = true
}
