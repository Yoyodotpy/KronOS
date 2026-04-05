package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
)

func ln(args []string) {
	symlink := flag.Bool("s", false, "create symlink instead of hard link")
	force := flag.Bool("f", false, "remove dest file")

	flag.Parse()

	dirs := flag.Args()

	if len(dirs) < 2 {
		err := errors.New("bro i need you to give me something")
		checkerr(err)
	}

	src := dirs[0]
	dest := dirs[1]

	info, err := os.Stat(dest)
	if err == nil && info.IsDir() {
		dest = filepath.Join(dest, filepath.Base(src))
	} else if os.IsNotExist(err) {
	} else {
		checkerr(err)
	}

	if *force {
		os.Remove(dest)
	}

	var linkerr error
	if *symlink {
		linkerr = os.Symlink(src, dest)
	} else {
		linkerr = os.Link(src, dest)
	}
	checkerr(linkerr)
}
