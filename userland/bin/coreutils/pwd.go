package main

import (
	"fmt"
	"os"
)

func pwd(_ []string) {
	dir, err := os.Getwd()
	fmt.Println(dir)
	checkerr(err)
}
