package main

import (
	_ "embed"
	"fmt"
	"image"
	"image/draw"
	"os"
	"time"

	"github.com/d21d3q/framebuffer"
	"github.com/nfnt/resize"
)

//go:embed anim/badapple.bin
var bafile []byte
var fb *framebuffer.Image

func main() {
	fb, _ = framebuffer.Open("/dev/fb0")
	defer fb.Close()

	for i := range len(bafile) / (640 * 480) {
		drawframe(i)
		time.Sleep(time.Millisecond * 50)
	}
}

func drawframe(i int) {
	w, h := 640, 480
	framesize := w * h
	start := i * framesize
	pixeldata := bafile[start : start+framesize]

	img := &image.Gray{
		Pix:    pixeldata,
		Stride: w,
		Rect:   image.Rect(0, 0, w, h),
	}

	target := resize.Resize(uint(fb.Bounds().Dx()), uint(fb.Bounds().Dy()), img, resize.NearestNeighbor)
	draw.Draw(fb, fb.Bounds(), target, image.Point{}, draw.Src)
}

func checkerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
