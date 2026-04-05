package main

import (
	"syscall"

	flag "github.com/spf13/pflag"
)

func mount(_ []string) {
	bind := flag.Bool("bind", false, "mount a subtree somewhere else")
	flags := flag.Int("f", 0, "what flag do you want to use, 1 for read only and 0 for read/write")
	fstype := "ext4"

	flag.Parse()

	args := flag.Args()

	drv := args[0]
	targ := args[1]

	data := ""

	mntflags := uintptr(*flags)
	if *bind {
		mntflags |= syscall.MS_BIND
		fstype = ""
	}

	err := syscall.Mount(drv, targ, fstype, uintptr(mntflags), data)
	checkerr(err)
}
