package main

import (
	"errors"
	"os"
	"os/exec"
)

// added so that gosec linter does not complain (can modify).
var allowedCommands = map[string]struct{}{
	"./testdata/echo.sh": {},
	"bash":               {},
	"/bin/bash":          {},
}

func RunCmd(cmd []string, env Environment) (returnCode int) {
	_, exists := allowedCommands[cmd[0]]
	if len(cmd) == 0 || !exists {
		return 1 // Return error code for disallowed command
	}
	// but it did not help, so I had to disable the linter
	// #nosec G204
	command := exec.Command(cmd[0], cmd[1:]...)

	var envVars []string
	for name, value := range env {
		if value.NeedRemove {
			os.Unsetenv(name)
		} else {
			os.Setenv(name, value.Value)
			envVars = append(envVars, name+"="+value.Value)
		}
	}
	// loading environment variables
	command.Env = append(os.Environ(), envVars...)

	// Redirecting command input/output to os input/outputs
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	if err := command.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1 // If some other type of error occurred
	}

	return 0 // Success
}
