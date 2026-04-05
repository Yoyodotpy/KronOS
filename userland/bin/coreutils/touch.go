package main

import (
	"errors"
	"os"
	"time"
)

func touch(args []string) {
	if len(args) >= 1 {
		target, err := os.OpenFile(args[0], os.O_RDONLY|os.O_CREATE, 0666)
		checkerr(err)
		target.Close()

		now := time.Now()
		err = os.Chtimes(os.Args[1], now, now)
		checkerr(err)
	} else {
		err := errors.New("you need to give me a name")
		checkerr(err)
	}
}
