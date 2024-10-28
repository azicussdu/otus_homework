package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		_, _ = fmt.Fprintf(os.Stderr, "No enough arguments\n")
		os.Exit(1)
	}

	envDir := os.Args[1]
	command := os.Args[2]
	args := os.Args[3:]

	env, err := ReadDir(envDir)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error loading environment variables: %v\n", err)
		os.Exit(1)
	}

	returnCode := RunCmd(append([]string{command}, args...), env)
	os.Exit(returnCode)
}
