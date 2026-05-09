package main

import (
	"image"
	"image/draw"
	"net/rpc"
	"os"
	"time"

	"github.com/d21d3q/framebuffer"
	"github.com/holoplot/go-evdev"
)

type fbuffer struct {
	fb     *framebuffer.Image
	db     draw.Image
	kbchan chan *evdev.InputEvent
}

var PID int
var client *rpc.Client

func main() {
	SetVirtualFileSystem(embeddedwad)

	var err error
	client, err = rpc.Dial("unix", "/tmp/wm.sock")
	checkerr(err)

	fb0, err := framebuffer.Open("/dev/fb0")
	checkerr(err)
	defer fb0.Close()

	kb0, err := evdev.Open("/dev/input/event1")
	checkerr(err)

	err = kb0.Grab()
	checkerr(err)

	defer kb0.Revoke()
	defer kb0.Close()

	eventChan := make(chan *evdev.InputEvent, 100)

	PID = os.Getpid()

	go func() {
		for {
			event, err := kb0.ReadOne()
			if err != nil {
				return
			}
			eventChan <- event
		}
	}()

	driver := &fbuffer{
		fb:     fb0,
		kbchan: eventChan,
		db:     image.NewRGBA(fb0.Bounds()),
	}

	var reply WindowReply
	defer client.Call("Winman.CloseWin", 1, &reply)

	Run(driver)
}

func (f *fbuffer) DrawFrame(img *image.RGBA) {
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
	time.Sleep(time.Microsecond * 100)
}

func (f *fbuffer) GetEvent(ev *DoomEvent) bool {
	select {
	case event := <-f.kbchan:
		if event.Type == evdev.EV_KEY {
			switch event.Value {
			case 0:
				ev.Type = Ev_keyup
			case 1:
				ev.Type = Ev_keydown
			default:
				return false
			}
			ev.Key = translate(event.Code)
			return true
		}
		return false
	default:
		return false
	}
}

func translate(key evdev.EvCode) uint8 {
	switch key {
	case evdev.KEY_W, evdev.KEY_UP:
		return KEY_UPARROW1
	case evdev.KEY_S, evdev.KEY_DOWN:
		return KEY_DOWNARROW1
	case evdev.KEY_A, evdev.KEY_LEFT:
		return KEY_LEFTARROW1
	case evdev.KEY_D, evdev.KEY_RIGHT:
		return KEY_RIGHTARROW1

	case evdev.KEY_LEFTCTRL, evdev.KEY_SPACE:
		return KEY_FIRE1
	case evdev.KEY_ENTER:
		return KEY_ENTER
	case evdev.KEY_ESC:
		return KEY_ESCAPE
	case evdev.KEY_E, evdev.KEY_RIGHTCTRL:
		return KEY_USE1
	case evdev.KEY_Y:
		return 'y'
	case evdev.KEY_N:
		return 'n'

	case evdev.KEY_1:
		return '1'
	case evdev.KEY_2:
		return '2'
	case evdev.KEY_3:
		return '3'
	case evdev.KEY_4:
		return '4'
	case evdev.KEY_5:
		return '5'
	case evdev.KEY_6:
		return '6'
	case evdev.KEY_7:
		return '7'
	case evdev.KEY_8:
		return '8'
	case evdev.KEY_9:
		return '9'
	case evdev.KEY_0:
		return '0'
	}
	return 0
}

func (f *fbuffer) CacheSound(name string, data []byte) {}

func (f *fbuffer) PlaySound(name string, channel, vol, sep int) {}

func (f *fbuffer) SetTitle(title string) {}

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}
