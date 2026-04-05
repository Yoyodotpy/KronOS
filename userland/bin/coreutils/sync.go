package main

import "syscall"

func sync(args []string) {
	syscall.Sync()
}
