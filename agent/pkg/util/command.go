package util

import (
	"bytes"
	"os/exec"
)

func ExecCommand(command string, args ...string) (string, error) {
	// Create an *exec.Cmd instance for the given command and its arguments
	cmd := exec.Command(command, args...)

	// A buffer to capture the standard output
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	// Run the command and wait for it to finish
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Return the captured standard output as a string
	return stdout.String(), nil
}
