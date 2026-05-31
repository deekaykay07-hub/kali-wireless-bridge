package hardware

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
)

// CommandResult represents output from an executed command
type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// ExecuteCommand runs a command and returns the result.
// In the future this will support streaming and privileged execution.
func ExecuteCommand(ctx context.Context, command string, args ...string) (*CommandResult, error) {
	cmd := exec.CommandContext(ctx, command, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var outBuf, errBuf strings.Builder

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			outBuf.WriteString(scanner.Text() + "\n")
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			errBuf.WriteString(scanner.Text() + "\n")
		}
	}()

	err = cmd.Wait()

	result := &CommandResult{
		Stdout:   outBuf.String(),
		Stderr:   errBuf.String(),
		ExitCode: 0,
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
	}

	return result, nil
}
