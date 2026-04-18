package main

import (
	"image"
	"image/color"
	"regexp"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var ansiRegex = regexp.MustCompile(`[\x1b\x9b][[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]`)

type FBWrapper struct{}

func (f *FBWrapper) ColorModel() color.Model { return color.RGBAModel }
func (f *FBWrapper) Bounds() image.Rectangle { return image.Rect(0, 0, 2000, 2000) }
func (f *FBWrapper) At(x, y int) color.Color {
	return color.RGBA{}
}
func (f *FBWrapper) Set(x, y int, c color.Color) {
	r, g, b, _ := c.RGBA()
	drawpix(x, y, int(r>>8), int(g>>8), int(b>>8))
}

var fbwrapper = &FBWrapper{}

func drawtext(text string, x int, y int, clr color.Color) {
	drawer := &font.Drawer{
		Dst:  fbwrapper,
		Src:  image.NewUniform(clr),
		Face: terminalface,
		Dot:  fixed.P(x, y+24),
	}
	drawer.DrawString(text)
}

func drawmultiline(lines []string, clr color.Color) {
	lineHeight := 30

	for i, line := range lines {
		drawtext(line, 10, i*lineHeight, clr)
	}
}

func cleanstring(str string) string {
	clean := ansiRegex.ReplaceAllString(str, "")
	clean = strings.ReplaceAll(clean, "\t", "    ")
	clean = strings.ReplaceAll(clean, "\r", "")
	return clean
}

func striplines(s string, maxlines int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > maxlines {
		return strings.Join(lines[len(lines)-maxlines:], "\n")
	}
	return s
}
