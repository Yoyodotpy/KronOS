package main

import (
	lua "github.com/yuin/gopher-lua"
)

func lua_cmd(args []string) {
	L := lua.NewState()
	defer L.Close()
	for _, file := range args {
		err := L.DoFile(file)
		checkerr(err)
	}
}
