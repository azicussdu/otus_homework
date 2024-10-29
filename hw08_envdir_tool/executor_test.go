package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	// A temporary env vars.
	env := Environment{
		"FOO":   {Value: "foo_value", NeedRemove: false},
		"BAR":   {Value: "bar_value", NeedRemove: false},
		"UNSET": {Value: "uns_value", NeedRemove: true}, // This should be removed.
	}

	// Define bash command to print env vars
	cmd := []string{"bash", "-c", "echo FOO=$FOO; echo BAR=$BAR; echo UNSET=$UNSET"}

	// because output of command is redirected to os.Stdout I should write it to buffer
	var buf bytes.Buffer
	multiWriter := io.MultiWriter(os.Stdout, &buf)
	originalStdout := os.Stdout

	// here need to create a temporary file to redirect output of command ti file first
	tempFile, err := os.CreateTemp("", "temp_stdout")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return
	}
	defer os.Remove(tempFile.Name())

	// Redirect output to the temporary file
	os.Stdout = tempFile

	// here I run the command, and output is generated and redirected to temp file
	exitCode := RunCmd(cmd, env)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	// Restore the original os.Stdout
	os.Stdout = originalStdout
	// read output from temp file
	_, err = tempFile.Seek(0, 0) // Go back to the start of the file
	if err != nil {
		fmt.Println("Error seeking temp file:", err)
		return
	}
	_, err = io.Copy(multiWriter, tempFile)
	if err != nil {
		fmt.Println("Error copying to multiWriter:", err)
		return
	}

	// Get the bytes from the buffer as a slice
	outputBytes := buf.Bytes()
	// Convert the bytes to a string
	output := string(outputBytes)

	// Check that the output contains the expected values.
	if !strings.Contains(output, "foo_value") {
		t.Error("expected FOO to be set correctly")
	}
	if !strings.Contains(output, "bar_value") {
		t.Error("expected BAR to be set correctly")
	}
	if strings.Contains(output, "uns_value") {
		t.Error("expected UNSET to be removed")
	}
}
