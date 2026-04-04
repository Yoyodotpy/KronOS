package main

import (
	"fmt"
	"os"
)

func main() {
	dir, err := os.Getwd()
	fmt.Println(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
