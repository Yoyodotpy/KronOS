package main

import (
	"image/color"
	"os"
	"strings"
	"sync"
)

var termstate string
var maxbuffersize int
var prevtermstate string

var termMu sync.Mutex

func initterm() {
	go func() {
		buf := make([]byte, 1)
		for {
			n, err := os.Stdin.Read(buf)
			checkerr(err)
			p.Write(buf[:n])
		}
	}()

	go func() {
		for {
			//time.Sleep(5 * time.Millisecond)
			b := make([]byte, 1024)
			n, err := p.Read(b)
			checkerr(err)
			out := cleanstring(string(b[:n]))

			termMu.Lock()
			termstate += out
			termstate = striplines(termstate, maxbuffersize)

			if strings.Contains(string(b), "\033[H\033[2J") {
				clear()
			}
			termMu.Unlock()
		}
	}()
}

func updateterm() bool {
	termMu.Lock()
	current := termstate
	termMu.Unlock()

	if current == prevtermstate {
		return false
	}

	drawrect(0, 0, int(scrinfo.Xres), int(scrinfo.Yres)-1, 0, 0, 0)
	lines := strings.Split(current, "\n")
	drawmultiline(lines, color.White)

	prevtermstate = current
	return true
}

func clear() {
	drawrect(0, 0, int(scrinfo.Xres), int(scrinfo.Yres)-1, 0, 0, 0)
	termstate = "> "
}
