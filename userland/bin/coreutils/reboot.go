package main

import (
	"syscall"
)

func reboot(_ []string) {
	syscall.Sync()
	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
	checkerr(err)
}
