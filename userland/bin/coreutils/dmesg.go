package main

import (
	"fmt"
	"os"
)

func dmesg(_ []string) {
	file, err := os.Open("/dev/kmsg")
	checkerr(err)
	defer file.Close()

	buf := make([]byte, 8192)
	for {
		text, err := file.Read(buf)
		checkerr(err)
		fmt.Println(string(buf[:text]))
	}
}
