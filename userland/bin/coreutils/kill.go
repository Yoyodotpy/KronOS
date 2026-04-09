package main

import (
	"errors"
	"os"
	"strconv"
)

func kill(args []string) {
	pid, err := strconv.Atoi(args[0])
	checkerr(err)

	if len(args) < 1 {
		err = errors.New("give me a pid (ls /proc)")
		checkerr(err)
	}

	process, err := os.FindProcess(pid)
	checkerr(err)

	err = process.Kill()
	checkerr(err)
}
