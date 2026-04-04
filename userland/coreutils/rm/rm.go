package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	var err []error
	if len(os.Args) > 1 {
		for _, i := range os.Args[1:] {
			e := os.Remove(i)
			if e != nil {
				err = append(err, e)
			}
		}
	} else {
		err = append(err, errors.New("path required"))
	}
	for _, i := range err {
		checkerr(i)
	}
}

func checkerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
