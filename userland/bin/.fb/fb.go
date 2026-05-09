package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/creack/pty"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/term"
)

//go:embed fonts/term.ttf
var fontBytes []byte

var p *os.File

func main() {
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	f, err := opentype.Parse(fontBytes)
	checkerr(err)
	terminalface, _ = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	fb, err := os.OpenFile(fbdev, os.O_RDWR, 0)
	checkerr(err)
	defer fb.Close()

	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,                 // The operation type
		fb.Fd(),                           // The "file" (/dev/fb0)
		0x4600,                            // FBIOGET_VSCREENINFO
		uintptr(unsafe.Pointer(&scrinfo)), // Pointer to our memory
	)
	if errno > 0 {
		checkerr(errno)
	}

	_, _, errno = syscall.Syscall(
		syscall.SYS_IOCTL,
		fb.Fd(),
		0x4602,
		uintptr(unsafe.Pointer(&fixinfo)),
	)
	if errno > 0 {
		checkerr(errno)
	}

	size := int(fixinfo.SMemLen)
	if size < 0 {
		err := errors.New("vro the screen is NOT smaller than 0")
		checkerr(err)
	}

	oldstate, err := term.MakeRaw(int(os.Stdin.Fd()))
	checkerr(err)
	defer term.Restore(int(os.Stdin.Fd()), oldstate)

	sh := exec.Command("/bin/sh")
	p, err = pty.Start(sh)
	checkerr(err)
	defer sh.Process.Kill()

	maxbuffersize = int(scrinfo.Yres / 30)

	screen, err = syscall.Mmap(int(fb.Fd()), 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	checkerr(err)
	defer syscall.Munmap(screen)

	frame = make([]byte, size)

	drawrect(0, 0, int(scrinfo.Xres), int(scrinfo.Yres)-1, 0, 0, 0)
	initterm()
	for {
		tick()
	}
}

func checkerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
