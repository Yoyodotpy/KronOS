package main

import (
	"fmt"
	"os"
)

func main() {
	dir, err := os.Getwd()
	checkerr(err)
	serveDir(dir)
}

func serveDir(dir string) {
	f, err := os.OpenFile(dir, os.O_RDONLY, 0666)
	checkerr(err)
	files, err := f.Readdirnames(0)
	checkerr(err)
	for _, file := range files {
		fmt.Print(" " + file)
	}
}

func checkerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
