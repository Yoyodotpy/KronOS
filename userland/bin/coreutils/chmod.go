package main

import (
	"errors"
	"os"
	"strconv"
)

func chmod(args []string) {
	if len(args) < 2 {
		err := errors.New("i need something to work with, please")
		checkerr(err)
	}

	mode, err := strconv.ParseUint(args[0], 8, 32)
	checkerr(err)

	fmode := os.FileMode(mode)

	err = os.Chmod(args[1], fmode)
	checkerr(err)
}
