package main

import (
	"fmt"
	"image"
	"image/draw"
	"io"
	"os"
	"time"

	"github.com/d21d3q/framebuffer"
	"github.com/fogleman/gg"
)

var screen *gg.Context

func main() {
	fb, _ := framebuffer.Open("/dev/fb0")
	defer fb.Close()

	screen = gg.NewContext(fb.Bounds().Dx(), fb.Bounds().Dy())

	go func() {
		_, err := io.Copy(vterm, os.Stdin)
		checkerr(err)
	}()

	for {

		screen.DrawCircle(100, 100, 400)
		screen.SetRGB(0, 100, 0)
		screen.Fill()
		draw.Draw(fb, fb.Bounds(), screen.Image(), image.Point{0, 0}, draw.Src)
		time.Sleep(time.Millisecond * 16)
	}
}

func checkerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
