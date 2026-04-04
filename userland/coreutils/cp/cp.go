package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	src := os.Args[1]
	dest := os.Args[2]

	file, err := os.Open(src)
	checkerr(err)
	defer file.Close()

	target, err := os.Create(dest)
	checkerr(err)
	defer target.Close()

	_, err = io.Copy(target, file)
	checkerr(err)
}

func checkerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
