package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var fbdev = "/dev/fb0"

type fbBitfield struct {
	Offset   uint32
	Length   uint32
	MsbRight uint32
}

type fbVarScreeninfo struct {
	Xres         uint32 // visible resolution
	Yres         uint32
	XresVirtual  uint32 // virtual resolution
	YresVirtual  uint32
	Xoffset      uint32 // offset from virtual to visible
	Yoffset      uint32
	BitsPerPixel uint32     // guess what this is!
	Grayscale    uint32     // 0 = color, 1 = grayscale, >1 = VIDEO_MODE_COLOR
	Red          fbBitfield // bitfield in fb mem if true color,
	Green        fbBitfield // else only length is significant
	Blue         fbBitfield
	Transp       fbBitfield
	Nonstd       uint32 // != 0 Non standard pixel format
	Activate     uint32 // see FB_ACTIVATE_*
	Height       uint32 // height of picture in mm
	Width        uint32 // width of picture in mm
	AccelFlags   uint32 // (obsolete) see fb_info.flags
	/* Timing: All values in pixclocks, except pixclock itself */
	Pixclock    uint32 // pixel clock in ps (pico seconds)
	LeftMargin  uint32 // time from sync to picture
	RightMargin uint32 // time from picture to sync
	UpperMargin uint32 // time from sync to picture
	LowerMargin uint32
	HsyncLen    uint32    // length of horizontal sync
	VsyncLen    uint32    // length of vertical sync
	Sync        uint32    // see FB_SYNC_*
	Vmode       uint32    // see FB_VMODE_*
	Rotate      uint32    // angle we rotate counter clockwise
	Colorspace  uint32    // colorspace for FOURCC-based modes
	Reserved    [4]uint32 // Reserved for future compatibility
}

type fbFixScreenInfo struct {
	ID           [16]byte // Screen ID, e.g. "TT Builtin"
	SMemStart    uintptr  // Frame buffer mem (physical address)
	SMemLen      uint32   // Frame buffer mem (physical address)
	Type         uint32   // Framebuffer type (see FB_TYPE_*)
	TypeAux      uint32   // "Interleave for interleaved Planes"
	Visual       uint32   // ??? (see FB_VISUAL_)
	XPanStep     uint16
	YPanStep     uint16
	YWrapStep    uint16
	LineLength   uint32  // Length of a line in bytes
	MmioStart    uintptr // Memory mapped I/O (physical address)
	MmioLen      uint32  // Memory mapped I/O (physical address)
	Accel        uint32  // "Indicate to driver which specific chip/card we have"
	Capabilities uint16  // contains filtered or unexported fields
}

var scrinfo fbVarScreeninfo
var fixinfo fbFixScreenInfo
var screen []byte

func main() {
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

	screen, err = syscall.Mmap(int(fb.Fd()), 0, size, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	checkerr(err)
	defer syscall.Munmap(screen)

	//drawpix(int(scrinfo.Xres)/2, int(scrinfo.Yres)/2, 255, 0, 0)
	//fmt.Println("Pixel drawn. Press Enter to exit...")
	//fmt.Scanln()

	//bpp := scrinfo.BitsPerPixel / 8
	//roffset := int(scrinfo.Red.Offset) / 8
	//groffset := int(scrinfo.Green.Offset) / 8
	//bloffset := int(scrinfo.Blue.Offset) / 8
	//c := 0
	//for i := 0; i < len(screen); i += int(bpp) {
	//	screen[i+bloffset] = 0            //blue
	//	screen[i+groffset] = 0 + byte(c)  //green
	//	screen[i+roffset] = 255 - byte(c) //red
	//	c++
	//}

	//centerX := int(scrinfo.Xres) / 2
	//centerY := int(scrinfo.Yres) / 2

	drawrect(1, 1, int(scrinfo.Xres)/2, int(scrinfo.Yres)-1, 255, 0, 0)
}

func checkerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getpos(x int, y int) int {
	pitch := int(fixinfo.LineLength)
	bpp := int(scrinfo.BitsPerPixel / 8)
	loc := (y * pitch) + (x * bpp)
	return loc
}

func drawpix(x int, y int, r int, g int, b int) {
	//bpp := scrinfo.BitsPerPixel / 8
	roffset := int(scrinfo.Red.Offset) / 8
	groffset := int(scrinfo.Green.Offset) / 8
	bloffset := int(scrinfo.Blue.Offset) / 8
	//troffset := int(scrinfo.Transp.Offset) / 8

	loc := getpos(x, y)
	screen[loc+roffset] = byte(r)
	screen[loc+groffset] = byte(g)
	screen[loc+bloffset] = byte(b)
}

func drawrect(x1 int, y1 int, x2 int, y2 int, r int, g int, b int) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			drawpix(x, y, r, g, b)
		}
	}
}
