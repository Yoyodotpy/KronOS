package main

import (
	"fmt"
	"os"
	"strconv"
)

func stat(args []string) {
	info, err := os.Stat(args[0])
	checkerr(err)

	fmt.Println("File: " + info.Name())
	fmt.Println("Size: " + strconv.FormatInt(info.Size(), 10))
	fmt.Println("Access: " + info.Mode().String())
	fmt.Println("Modify: " + info.ModTime().String())
	fmt.Println("Dir: " + strconv.FormatBool(info.IsDir()))

}
