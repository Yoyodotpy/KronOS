package main

import (
	_ "embed"
	"fmt"
	"image"
	"image/draw"
	"net/rpc"
	"os"
	"time"
)

//go:embed badapple.bin
var bafile []byte

var PID int
var client *rpc.Client

func main() {
	var err error
	client, err = rpc.Dial("unix", "/tmp/wm.sock")
	checkerr(err)

	PID = os.Getpid()

	var reply WindowReply
	title := Title{
		Title: "badapple do dodoo odo dodo oodododo doo",
		ID:    PID,
	}
	err = client.Call("Winman.SetTitle", title, &reply)
	checkerr(err)

	for i := range len(bafile) / (144 * 256) {
		drawframe(i)
		time.Sleep(time.Millisecond * 100)
	}
}

func drawframe(i int) {
	w, h := 256, 144
	framesize := w * h
	start := i * framesize
	pixeldata := bafile[start : start+framesize]

	gray := &image.Gray{
		Pix:    pixeldata,
		Stride: w,
		Rect:   image.Rect(0, 0, w, h),
	}

	img := image.NewRGBA(gray.Bounds())
	draw.Draw(img, img.Bounds(), gray, gray.Bounds().Min, draw.Src)

	info := &Winfo{
		ID:     PID,
		Width:  img.Bounds().Dx(),
		Height: img.Bounds().Dy(),
		Stride: img.Stride,
		Pixels: img.Pix,
	}

	var reply WindowReply
	err := client.Call("Winman.CreateWin", info, &reply)
	checkerr(err)
}

func checkerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
