package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open(os.Args[1])
	checkerr(err)

	text, err := io.ReadAll(file)
	checkerr(err)

	fmt.Println(string(text))
}

func checkerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
