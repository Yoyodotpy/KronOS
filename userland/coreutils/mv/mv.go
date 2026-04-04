package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	src := os.Args[1]
	dest := os.Args[2]

	info, err := os.Stat(dest)
	if err == nil && info.IsDir() {
		dest = filepath.Join(dest, filepath.Base(src))
	} else if os.IsNotExist(err) {
	} else {
		checkerr(err)
	}

	err = os.Rename(src, dest)
	checkerr(err)
}

func checkerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
