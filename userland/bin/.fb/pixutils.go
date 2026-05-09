package main

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
	frame[loc+roffset] = byte(r)
	frame[loc+groffset] = byte(g)
	frame[loc+bloffset] = byte(b)
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
			if x < int(scrinfo.Xres) && y < int(scrinfo.Yres) {
				drawpix(x, y, r, g, b)
			}
		}
	}
}
