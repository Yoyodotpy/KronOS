package main

import (
	"errors"
	"os"
)

func mkdir(args []string) {
	if len(args) >= 1 {
		os.Mkdir(args[0], 0755)
	} else {
		err := errors.New("i need something")
		checkerr(err)
	}
}
