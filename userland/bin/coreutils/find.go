package main

import (
	"os"
)

var files []string

func find(args []string) {
	dir, err := os.Getwd()
	if len(args) >= 1 {
		dir = args[0]
	}
	checkerr(err)
	files = listdir(dir)
	//loopdir()
}

func listdir(dir string) []string {
	f, err := os.OpenFile(dir, os.O_RDONLY, 0666)
	checkerr(err)
	files, err := f.Readdirnames(0)
	checkerr(err)
	return files
}

func loopdir(folders []string) []string {
	var files []string

	for _, file := range folders {
		info, err := os.Stat(file)
		checkerr(err)
		if info.IsDir() {
			branch := listdir(file)
			loopdir(branch)
		} else {
			files = append(files, file)
		}

	}
	return files
}
