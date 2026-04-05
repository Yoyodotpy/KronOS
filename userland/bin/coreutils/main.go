package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type cmdfunc func(args []string)

var cmds = make(map[string]cmdfunc)

var utildir = "/bin/coreutils"

func init() {
	cmds["coreutils"] = coreutils
	cmds["cp"] = cp
	cmds["echo"] = echo
	cmds["ln"] = ln
	cmds["ls"] = ls
	cmds["mkdir"] = mkdir
	cmds["mv"] = mv
	cmds["cat"] = cat
	cmds["pwd"] = pwd
	cmds["reboot"] = reboot
	cmds["rm"] = rm
	cmds["shutdown"] = shutdown
	cmds["touch"] = touch
	cmds["mount"] = mount
	cmds["umount"] = umount
	cmds["sync"] = sync
}

func main() {
	cmd := filepath.Base(os.Args[0])

	if funct, printer := cmds[cmd]; printer {
		funct(os.Args[1:])
	} else {
		fmt.Fprintf(os.Stderr, "could not find %s\n", cmd)
		os.Exit(1)
	}
}

func coreutils(args []string) {
	bin := "/bin"

	for name := range cmds {
		os.Symlink(utildir, filepath.Join(bin, name))
	}
}

func checkerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
