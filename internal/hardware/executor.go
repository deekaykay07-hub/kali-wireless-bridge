package hardware

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// CommandResult represents output from an executed command
type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// ExecuteCommand runs a command and returns the result.
// It automatically chooses the right shell based on OS.
func ExecuteCommand(ctx context.Context, command string, args ...string) (*CommandResult, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// On Windows, use PowerShell for better compatibility
		psArgs := append([]string{"-NoProfile", "-NonInteractive", "-Command"}, command)
		psArgs = append(psArgs, args...)
		cmd = exec.CommandContext(ctx, "powershell.exe", psArgs...)
	} else {
		cmd = exec.CommandContext(ctx, command, args...)
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var outBuf, errBuf strings.Builder

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			outBuf.WriteString(scanner.Text() + "\n")
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
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
		} else {
			result.ExitCode = 1
		}
	}

	return result, nil
}

// ListWirelessInterfaces returns available wireless interfaces
func ListWirelessInterfaces() ([]string, error) {
	if runtime.GOOS == "windows" {
		// Windows implementation using netsh
		cmd := exec.Command("netsh", "wlan", "show", "interfaces")
		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		// Very basic parsing - we can improve this
		lines := strings.Split(string(output), "\n")
		var interfaces []string
		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), "name") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					interfaces = append(interfaces, strings.TrimSpace(parts[1]))
				}
			}
		}
		return interfaces, nil
	}

	// Linux / macOS fallback
	cmd := exec.Command("iw", "dev")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list interfaces: %w", err)
	}

	// Simple parsing for Linux
	var interfaces []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Interface ") {
			interfaces = append(interfaces, strings.TrimPrefix(line, "Interface "))
		}
	}
	return interfaces, nil
}
