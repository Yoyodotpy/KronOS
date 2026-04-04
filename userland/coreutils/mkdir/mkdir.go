package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		os.Mkdir(os.Args[1], 0755)
	} else {
		fmt.Println("please give me something")
		os.Exit(1)
	}
}
