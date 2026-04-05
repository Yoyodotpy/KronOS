package main

import (
	"io"
	"os"
)

func cp(args []string) {
	src := args[0]
	dest := args[1]

	file, err := os.Open(src)
	checkerr(err)
	defer file.Close()

	target, err := os.Create(dest)
	checkerr(err)
	defer target.Close()

	_, err = io.Copy(target, file)
	checkerr(err)
}
