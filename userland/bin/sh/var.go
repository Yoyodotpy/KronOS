package main

import (
	"os"
	"strings"
)

func getval(arg string) string {
	if strings.HasPrefix(arg, "$(") && strings.HasSuffix(arg, ")") {
		cmd := arg[2 : len(arg)-1]
		output, _ := execInput(cmd, true)
		return strings.TrimSuffix(output, "\n")
	} else if name, found := strings.CutPrefix(arg, "$"); found {
		return os.Getenv(name)
	} else {
		return arg
	}
}
