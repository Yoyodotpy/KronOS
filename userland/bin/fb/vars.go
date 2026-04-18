package main

import "golang.org/x/image/font"

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
var frame []byte
var terminalface font.Face
