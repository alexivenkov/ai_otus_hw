package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	c := buildCommand(cmd[0], cmd[0:], env)

	if err := c.Run(); err != nil {
		var exitError *exec.ExitError

		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
	}

	return 0
}

func buildCommand(command string, args []string, env Environment) *exec.Cmd {
	cmd := exec.Command(command)

	for name, e := range env {
		if e.NeedRemove {
			os.Unsetenv(name)
			continue
		}
		os.Setenv(name, e.Value)
	}

	cmd.Env = os.Environ()

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	cmd.Args = args

	return cmd
}
