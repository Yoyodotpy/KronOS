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
}

type WindowReply struct {
	Success bool
	Focused bool
}

type Title struct {
	Title string
	ID    int
}

type window struct {
	ID      int
	X       int
	Y       int
	width   int
	height  int
	img     image.RGBA
	title   string
	focused bool
}

func (win *window) close() {
	proc, _ := os.FindProcess(win.ID)
	proc.Kill()
	delete(windows, win.ID)
	windowsneedrefresh = true
}

func (win *window) focuse() {
	for _, w := range windows {
		w.focused = false
	}
	win.focused = true
}
