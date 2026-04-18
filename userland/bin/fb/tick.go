package main

import "time"

func tick() {
	//drawtext("hello world", 0, 0, color.White)
	if updateterm() {
		copy(screen, frame)
	}
	time.Sleep(time.Millisecond * 16)
}
