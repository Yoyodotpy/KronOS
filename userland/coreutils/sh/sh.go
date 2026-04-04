package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		//read keyboard
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		_, err = execInput(input, false)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func execInput(input string, vared bool) (string, error) {
	input = strings.TrimSuffix(input, "\n")
	args := strings.Split(input, " ")
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
	} else {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		return "", cmd.Run()
	}
}
