package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	tcell "github.com/gdamore/tcell/v3"
	color "github.com/gdamore/tcell/v3/color"
)

type location struct {
	X int
	Y int
}

var curfile []string
var defstyle tcell.Style
var barstyle tcell.Style
var screen tcell.Screen
var cursor = location{X: 0, Y: 0}
var offset = location{X: 0, Y: 0}
var status = ""

func edit(args []string) {
	var err error

	if len(args) > 0 {
		filename := args[0]
		file, err := os.Open(filename)

		if err != nil {
			if os.IsNotExist(err) {
				file, err = os.Create(filename)
				checkerr(err)
			} else {
				checkerr(err)
			}
		}
		defer file.Close()

		text, err := io.ReadAll(file)
		checkerr(err)
		content := string(text)
		if content == "" {
			curfile = []string{""}
		} else {
			curfile = strings.Split(content, "\n")
		}
	} else {
		checkerr(errors.New("need filename"))
	}

	screen, err = tcell.NewScreen()
	checkerr(err)

	screen.Init()
	defer screen.Fini()
	screen.Sync()

	defstyle = tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	barstyle = tcell.StyleDefault.Background(color.White).Foreground(color.Black)
	screen.SetStyle(defstyle)

	curline := curfile[0]

	for {
		width, height := screen.Size()
		viewheight := height - 1
		viewwidth := width - 5

		if cursor.Y < offset.Y {
			offset.Y = cursor.Y
		} else if cursor.Y >= offset.Y+viewheight {
			offset.Y = cursor.Y - viewheight + 1
		}
		if cursor.X < offset.X {
			offset.X = cursor.X
		} else if cursor.X >= offset.X+viewwidth {
			offset.X = cursor.X - viewwidth + 1
		}

		screen.Clear()
		curfile[cursor.Y] = curline
		drawparag(curfile)
		drawcursor()
		drawbar()
		screen.Show()

		ev := <-screen.EventQ()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlC, tcell.KeyCtrlQ, tcell.KeyCtrlX:
				return
			case tcell.KeyCtrlS, tcell.KeyCtrlO:
				status = "saved"
				save(args[0])

			case tcell.KeyEnter:
				before := curline[:cursor.X]
				after := curline[cursor.X:]

				curfile[cursor.Y] = before
				curfile = append(curfile[:cursor.Y+1], append([]string{after}, curfile[cursor.Y+1:]...)...)

				cursor.Y += 1
				cursor.X = 0

				curline = curfile[cursor.Y]

			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if cursor.X > 0 {
					curline = curline[:cursor.X-1] + curline[cursor.X:]
					cursor.X -= 1
				} else if cursor.Y > 0 {
					prevLineLen := len(curfile[cursor.Y-1])
					curfile[cursor.Y-1] += curline

					curfile = append(curfile[:cursor.Y], curfile[cursor.Y+1:]...)

					cursor.Y -= 1
					cursor.X = prevLineLen
					curline = curfile[cursor.Y]
				}

			case tcell.KeyLeft:
				if cursor.X > 0 {
					cursor.X -= 1
				} else if cursor.Y > 0 {
					cursor.Y -= 1
					curline = curfile[cursor.Y]
					cursor.X = len(curline)
				}
			case tcell.KeyRight:
				if cursor.X < len(curline) {
					cursor.X += 1
				} else if cursor.Y < len(curfile)-1 {
					cursor.Y += 1
					curline = curfile[cursor.Y]
					cursor.X = 0
				}
			case tcell.KeyUp:
				if cursor.Y > 0 {
					curfile[cursor.Y] = curline
					cursor.Y -= 1
					curline = curfile[cursor.Y]
					if cursor.X > len(curline) {
						cursor.X = len(curline)
					}
				}
			case tcell.KeyDown:
				if cursor.Y < len(curfile)-1 {
					curfile[cursor.Y] = curline
					cursor.Y += 1
					curline = curfile[cursor.Y]
					if cursor.X > len(curline) {
						cursor.X = len(curline)
					}
				}

			default:
				if ev.Str() != "" {
					modable := []rune(curline)

					var before []rune
					var after []rune
					if cursor.X >= len(modable) {
						before = modable[:cursor.X]
						after = modable[cursor.X:]
					} else {
						before = modable[:cursor.X]
						after = modable[cursor.X+1:]
					}

					newline := string(before) + ev.Str() + string(after)
					curline = newline

					cursor.X += 1
				}
			}
		case *tcell.EventResize:
			screen.Sync()
		}
	}
}

func drawline(text string, line int) {
	for i, char := range text {
		screen.Put(i, line, string(char), defstyle)
	}
}

func drawparag(text []string) {
	width, height := screen.Size()
	viewheight := height - 1

	for i := range viewheight {
		lineindex := i + offset.Y
		if lineindex < len(text) {
			line := text[lineindex]

			if offset.X < len(line) {
				line = line[offset.X:]
			} else {
				line = ""
			}
			disp := fmt.Sprintf("%3d %s", lineindex, line)

			if len(disp) > width {
				disp = disp[:width]
			}

			drawline(disp, i)
		}
	}
}

func drawcursor() {
	screenx := (cursor.X - offset.X) + 4
	screeny := cursor.Y - offset.Y
	_, height := screen.Size()
	if screeny >= 0 && screeny < height-1 {
		char, _, _ := screen.Get(screenx, screeny)
		screen.Put(screenx, screeny, char, barstyle)
	}
}

func drawbar() {
	_, height := screen.Size()
	for i, a := range status {
		screen.Put(i, height-1, string(a), barstyle)
	}
	status = ""
}

func save(file string) {
	os.Create(file)
	data := strings.Join(curfile, "\n")
	err := os.WriteFile(file, []byte(data), 0644)
	checkerr(err)
}
