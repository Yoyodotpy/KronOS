package main

import (
	"errors"
	"syscall"

	flag "github.com/spf13/pflag"
)

func umount(_ []string) {
	force := flag.BoolP("force", "f", false, "force unmount")
	lazy := flag.BoolP("lazy", "l", false, "lazy unmount")

	flag.Parse()

	var flags int = 0

	if *force {
		flags |= syscall.MNT_FORCE
	}
	if *lazy {
		flags |= syscall.MNT_DETACH
	}

	args := flag.Args()

	if len(args) < 1 {
		err := errors.New("i need something")
		checkerr(err)
	}

	targ := args[0]

	err := syscall.Unmount(targ, flags)
	checkerr(err)
}
