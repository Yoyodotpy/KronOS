package main

import (
	"fmt"
	"io"
	"os"
)

func cat(args []string) {
	file, err := os.Open(args[0])
	checkerr(err)

	text, err := io.ReadAll(file)
	checkerr(err)

	fmt.Println(string(text))
}
