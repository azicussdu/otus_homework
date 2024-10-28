package main

import (
	"os/exec"
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

	exitCode := RunCmd(cmd, env)
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		t.Fatalf("error executing command: %v", err)
	}
	output := string(out)

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
