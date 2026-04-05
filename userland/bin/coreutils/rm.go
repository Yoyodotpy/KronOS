package main

import (
	"errors"
	"os"
)

func rm(args []string) {
	var err []error
	if len(args) >= 1 {
		for _, i := range args {
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
