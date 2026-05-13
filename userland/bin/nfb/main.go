package main

import (
	"embed"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"time"

	"github.com/d21d3q/framebuffer"
	"github.com/fogleman/gg"
	"github.com/holoplot/go-evdev"
)

//go:embed assets
var assets embed.FS
var screenbuffer *gg.Context

func main() {
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	defer fmt.Print("\033[H\033[2J")

	fb, err := framebuffer.Open("/dev/fb0")
	checkerr(err)
	defer fb.Close()

	screenbuffer = gg.NewContext(fb.Bounds().Dx(), fb.Bounds().Dy())

	mouseDev, err := evdev.Open("/dev/input/event2")
	checkerr(err)

	file, err := assets.Open("assets/cursor.png")
	checkerr(err)
	cursorsprite, _, err = image.Decode(file)
	checkerr(err)

	os.Remove("/tmp/wm.sock")
	rpc.Register(new(Winman))
	listener, err := net.Listen("unix", "/tmp/wm.sock")
	checkerr(err)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()

	go spawnclient("/bin/doom")
	go spawnclient("/bin/badapple")
	go mousemove(mouseDev, fb.Bounds().Dx(), fb.Bounds().Dy())

	screenbuffer.SetRGB(0.1, 0.1, 0.1)
	screenbuffer.Clear()
	draw.Draw(fb, fb.Bounds(), screenbuffer.Image(), image.Point{0, 0}, draw.Src)

	var lastx, lasty int
	ticker := time.NewTicker(16 * time.Millisecond)

	for range ticker.C {
		interactwindow()

		redraw := (lastx != cursor.x || lasty != cursor.y) || windowsneedrefresh

		if redraw {
			drawwindows(fb)
			windowsneedrefresh = false
			lastx, lasty = cursor.x, cursor.y
		}
	}
}

func checkerr(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func spawnclient(path string) {
	cmd := exec.Command(path)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	checkerr(err)
}
