package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		var builder strings.Builder
		for _, arg := range os.Args[1:] {
			builder.WriteString(arg + " ")
		}
		echoed := builder.String()
		fmt.Println(echoed)
	}
}
