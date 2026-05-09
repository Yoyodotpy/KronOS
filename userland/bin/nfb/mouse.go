package main

import (
	"github.com/holoplot/go-evdev"
)

type MouseState struct {
	x, y     int
	lmb, rmb bool
}

var cursor MouseState

func mousemove(dev *evdev.InputDevice, width, height int) {
	for {
		events, err := dev.ReadSlice(64)
		if err != nil {
			continue
		}
		for _, event := range events {
			if event.Type == evdev.EV_REL {
				switch event.Code {
				case evdev.REL_X:
					cursor.x += int(event.Value)
				case evdev.REL_Y:
					cursor.y += int(event.Value)
				}
			}
			if event.Type == evdev.EV_KEY {
				switch event.Code {
				case evdev.BTN_LEFT:
					cursor.lmb = event.Value != 0
				case evdev.BTN_RIGHT:
					cursor.rmb = event.Value != 0
				}
			}
		}

		if cursor.x < 0 {
			cursor.x = 0
		}
		if cursor.y < 0 {
			cursor.y = 0
		}
		if cursor.x > width {
			cursor.x = width
		}
		if cursor.y > height {
			cursor.y = height
		}
	}
}
