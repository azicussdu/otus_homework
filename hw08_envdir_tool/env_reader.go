package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue // skip if directory
		}

		name := file.Name()
		if strings.Contains(name, "=") {
			return nil, nil // ignore files with "=" in the name
		}

		f, err := os.Open(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		defer f.Close()

		reader := bufio.NewReader(f)
		line, err := reader.ReadString('\n')  // Read until \n
		line = strings.TrimSuffix(line, "\n") // Remove if you find newline

		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err // can't read file
		}

		if errors.Is(err, io.EOF) && line == "" { // if file is empty
			env[name] = EnvValue{NeedRemove: true}
			continue
		}

		value := strings.TrimRight(line, " \t")         // remove right side spaces and tabs
		value = strings.ReplaceAll(value, "\x00", "\n") // replace \0 to \n
		env[name] = EnvValue{Value: value, NeedRemove: false}
	}
	return env, nil
}
