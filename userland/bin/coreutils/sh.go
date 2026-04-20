package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func sh(args []string) {
	if len(args) > 0 {
		lua_cmd(args)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		//read keyboard
		input, err := reader.ReadString('\n')
		checkerr(err)

		_, err = execInput(input, false)
		checkerr(err)
	}
}

func execInput(input string, vared bool) (string, error) {
	input = strings.TrimSuffix(input, "\n")
	args := strings.Split(input, " ")
	var stdout string

	for i, arg := range args {
		if arg == ">" {
			stdout = args[i+1]
			args = append(args[:i], args[i+2:]...)
			break
		}
	}

	for i := range args {
		args[i] = getval(args[i])
	}

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return "", errors.New("path required")
		}
		return "", os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	case "clear":
		fmt.Print("\033[H\033[2J")
		return "", nil
	}

	cmd := exec.Command(args[0], args[1:]...)
	var output []byte
	var err error

	if vared {
		output, err = cmd.CombinedOutput()
		return string(output), err
	}

	if stdout != "" {
		f, err := os.Create(stdout)
		checkerr(err)
		defer f.Close()
		cmd.Stdout = f
	} else {
		cmd.Stdout = os.Stdout
	}

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return "", cmd.Run()
}
