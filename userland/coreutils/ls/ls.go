package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	var dir string
	dir, err := os.Getwd()
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	checkerr(err)
	serveDir(dir)
}

func serveDir(dir string) {
	f, err := os.OpenFile(dir, os.O_RDONLY, 0666)
	checkerr(err)
	files, err := f.Readdirnames(0)
	checkerr(err)
	var builder strings.Builder
	for _, file := range files {
		builder.WriteString(file + " ")
	}
	dirlist := builder.String()
	fmt.Println(dirlist)
}

func checkerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
