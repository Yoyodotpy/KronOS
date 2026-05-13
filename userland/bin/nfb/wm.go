package main

import (
	"image"
	"image/draw"
)

type Winman struct{}

var windows = make(map[int]*window)
var deadwindows = make(map[int]bool)
var windowsneedrefresh bool
var cursorsprite image.Image

var activewindow int = -1
var dragoffsetx, dragoffsety int

const barheight = 15

func (w *Winman) CreateWin(info *Winfo, reply *WindowReply) error {
	if deadwindows[info.ID] {
		reply.Success = false
		return nil
	}

	if win, exists := windows[info.ID]; exists {
		if win.img.Bounds().Dx() != info.Width || win.img.Bounds().Dy() != info.Height || win.img.Stride != info.Stride {
			rect := image.Rect(0, 0, info.Width, info.Height)
			tempimg := image.RGBA{
				Pix:    info.Pixels,
				Stride: info.Stride,
				Rect:   rect,
			}
			win.img = tempimg
		} else {
			copy(win.img.Pix, info.Pixels)
		}
	} else {
		rect := image.Rect(0, 0, info.Width, info.Height)
		windows[info.ID] = &window{
			ID: info.ID,
			X:  100,
			Y:  100,
			img: image.RGBA{
				Pix:    info.Pixels,
				Stride: info.Stride,
				Rect:   rect,
			},
			focused: true,
		}
	}
	windowsneedrefresh = true
	reply.Success = true
	reply.Focused = windows[info.ID].focused
	return nil
}

func (w *Winman) CloseWin(ID int, reply *bool) error {
	windows[ID].close()
	return nil
}

func (w *Winman) SetTitle(title Title, reply *WindowReply) error {
	_, exists := windows[title.ID]
	if !exists {
		windows[title.ID] = &window{
			ID:    title.ID,
			X:     100,
			Y:     100,
			title: title.Title,
			img: image.RGBA{
				Pix:    make([]byte, 10*10*4),
				Stride: 10 * 4,
				Rect:   image.Rect(0, 0, 10, 10),
			},
		}
	}
	windows[title.ID].title = title.Title
	windowsneedrefresh = true
	reply.Success = true
	return nil
}

func drawwindows(fb draw.Image) {
	screenbuffer.SetRGB(0.1, 0.1, 0.1)
	screenbuffer.Clear()

	var curwin window

	for _, window := range windows {
		if !window.focused {
			screenbuffer.SetRGB(0.3, 0.3, 0.3)
			screenbuffer.DrawRectangle(float64(window.X), float64(window.Y-barheight), float64(window.img.Bounds().Dx()), float64(barheight))
			screenbuffer.Fill()
			screenbuffer.SetRGB(0.7, 0, 0)
			screenbuffer.DrawRectangle(float64(window.X), float64(window.Y-barheight), float64(barheight), float64(barheight))
			screenbuffer.Fill()
			screenbuffer.SetRGB(1, 1, 1)
			screenbuffer.DrawString(window.title, float64(window.X+barheight+1), float64(window.Y-1))
			screenbuffer.Fill()
			screenbuffer.SetRGB(0.1, 0.1, 0.1)

			var frame image.Image = &window.img
			screenbuffer.DrawImageAnchored(frame, window.X, window.Y, 0, 0)
		} else {
			curwin = *window
		}
	}

	screenbuffer.SetRGB(0, 0.7, 0)
	screenbuffer.DrawRectangle(float64(curwin.X), float64(curwin.Y-barheight), float64(curwin.img.Bounds().Dx()), float64(barheight))
	screenbuffer.Fill()
	screenbuffer.SetRGB(0.7, 0, 0)
	screenbuffer.DrawRectangle(float64(curwin.X), float64(curwin.Y-barheight), float64(barheight), float64(barheight))
	screenbuffer.Fill()
	screenbuffer.SetRGB(1, 1, 1)
	screenbuffer.DrawString(curwin.title, float64(curwin.X+barheight+1), float64(curwin.Y-1))
	screenbuffer.Fill()
	screenbuffer.SetRGB(0.1, 0.1, 0.1)
	var frame image.Image = &curwin.img
	screenbuffer.DrawImageAnchored(frame, curwin.X, curwin.Y, 0, 0)

	screenbuffer.DrawImageAnchored(cursorsprite, cursor.x, cursor.y, 0, 0)

	windowframe := screenbuffer.Image()
	draw.Draw(fb, windowframe.Bounds(), windowframe, windowframe.Bounds().Min, draw.Src)
}

func interactwindow() {
	if cursor.lmb {
		if activewindow == -1 {
			var x bool
			id, x := getwindow(cursor.x, cursor.y)
			if id != -1 {
				activewindow = id
				win := windows[activewindow]
				win.focuse()
				dragoffsetx = cursor.x - win.X
				dragoffsety = cursor.y - win.Y
				if x {
					windows[activewindow].close()
					deadwindows[activewindow] = true
					activewindow = -1
				}
			}
		} else {
			win := windows[activewindow]
			win.X = cursor.x - dragoffsetx
			win.Y = cursor.y - dragoffsety
			windows[activewindow] = win
			windowsneedrefresh = true
		}
	} else {
		activewindow = -1
	}
}

func getwindow(x int, y int) (int, bool) {
	for id, win := range windows {
		if x >= win.X && x <= win.X+win.img.Bounds().Dx() && y >= win.Y-barheight && y <= win.Y {
			if x < (win.X + barheight) {
				return id, true
			}
			return id, false
		}
	}
	return -1, false
}
