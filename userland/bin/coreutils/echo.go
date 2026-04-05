package main

import (
	"fmt"
	"strings"
)

func echo(args []string) {
	if len(args) >= 1 {
		var builder strings.Builder
		for _, arg := range args {
			builder.WriteString(arg + " ")
		}
		echoed := builder.String()
		fmt.Println(echoed)
	}
}
