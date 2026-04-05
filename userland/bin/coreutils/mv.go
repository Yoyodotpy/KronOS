package main

import (
	"errors"
	"os"
	"path/filepath"
)

func mv(args []string) {
	if len(args) < 2 {
		err := errors.New("bro i need you to give me something")
		checkerr(err)
	}

	src := args[0]
	dest := args[1]

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
