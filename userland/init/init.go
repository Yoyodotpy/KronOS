package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var shell = "/bin/sh"

func main() {
	//fmt.Print("\033[H\033[2J")
	fmt.Println("i am alive")

	// Mount file system
	os.MkdirAll("/proc", 0755)
	syscall.Mount("proc", "/proc", "proc", 0, "")

	os.MkdirAll("/sys", 0755)
	syscall.Mount("sysfs", "/sys", "sysfs", 0, "")

	os.MkdirAll("/dev", 0755)
	syscall.Mount("devtmpfs", "/dev", "devtmpfs", 0, "")

	os.Setenv("PATH", "/bin:/sbin")

	for {
		execInput(shell)
	}
}

func execInput(input string) error {
	input = strings.TrimSuffix(input, "\n")
	args := strings.Split(input, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func checkerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
