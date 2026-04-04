package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args[1]) > 1 {
		target, err := os.OpenFile(os.Args[1], os.O_RDONLY|os.O_CREATE, 0666)
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

func checkerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
