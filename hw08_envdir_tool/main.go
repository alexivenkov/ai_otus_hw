package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 3 {
		fmt.Println("env path or executable command was not specify")
		os.Exit(0)
	}

	path := args[1]
	command := args[2:]

	env, err := ReadDir(path)
	if err != nil {
		fmt.Println("unable to read env dir")
		os.Exit(0)
	}

	exitCode := RunCmd(command, env)
	os.Exit(exitCode)
}
