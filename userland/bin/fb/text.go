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
var fontletters []*image.RGBA

func drawtext(text string, x int, y int, clr color.Color) {
	drawer := &font.Drawer{
		Src:  image.NewUniform(clr),
		Face: terminalface,
	}

	advance := drawer.MeasureString(text)
	width := advance.Ceil()
	height := 32

	if width <= 0 {
		return
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	drawer.Dst = img

	drawer.Dot = fixed.P(0, 24)
	drawer.DrawString(text)

	for ty := range height {
		for tx := range width {
			r, g, b, a := img.At(tx, ty).RGBA()
			if a > 0 {
				drawpix(x+tx, y+ty, int(r>>8), int(g>>8), int(b>>8))
			}
		}
	}
}

func drawmultiline(lines []string, clr color.Color) {
	lineHeight := 30

	for i, line := range lines {
		drawtext(line, 10, i*lineHeight, clr)
	}
}

func cleanstring(str string) string {
	clean := ansiRegex.ReplaceAllString(str, "")
	clean = strings.ReplaceAll(clean, "\t", "   ")
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

func loadfont(clr color.Color) {
	var let = "a"

	drawer := &font.Drawer{
		Src:  image.NewUniform(clr),
		Face: terminalface,
	}

	advance := drawer.MeasureString(let)
	width := advance.Ceil()
	height := 32

	if width <= 0 {
		return
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	drawer.Dst = img

	drawer.Dot = fixed.P(0, 24)
	drawer.DrawString(let)

	fontletters = append(fontletters, img)
}
