package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var shell = "/bin/sh"
var coreutils = "/bin/coreutils"

func main() {
	fmt.Print("\033[H\033[2J")
	fmt.Println("i am alive")

	// Mount file system
	os.MkdirAll("/proc", 0755)
	syscall.Mount("proc", "/proc", "proc", 0, "")

	os.MkdirAll("/sys", 0755)
	syscall.Mount("sysfs", "/sys", "sysfs", 0, "")

	os.MkdirAll("/dev", 0755)
	syscall.Mount("devtmpfs", "/dev", "devtmpfs", 0, "")

	os.MkdirAll("/dev/pts", 0755)
	syscall.Mount("devpts", "/dev/pts", "devpts", 0, "newinstance,ptmxmode=0666,mode=620")
	os.Symlink("/dev/pts/ptmx", "/dev/ptmx")

	os.Setenv("PATH", "/bin:/sbin")

	execInput("/bin/coreutils")

	syscall.Syscall(syscall.SYS_SYSLOG, 8, 0, uintptr(1))

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
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid:  true,
		Setctty: true,
		Ctty:    0,
	}
	return cmd.Run()
}

func checkerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
