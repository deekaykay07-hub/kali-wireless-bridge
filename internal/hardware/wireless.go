package hardware

import (
	"fmt"
	"os/exec"
	"strings"
)

// Interface represents a wireless network interface
type Interface struct {
	Name       string
	Type       string // managed, monitor, etc.
	MAC        string
	Supported  bool
	CanInject  bool
}

// ListWirelessInterfaces returns available wireless interfaces on the system
func ListWirelessInterfaces() ([]Interface, error) {
	// This is a basic implementation using `iw`
	// In production this should be more robust and cross-platform

	cmd := exec.Command("iw", "dev")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run 'iw dev': %w (is iw installed?)", err)
	}

	var interfaces []Interface
	lines := strings.Split(string(output), "\n")

	var current Interface
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Interface ") {
			if current.Name != "" {
				interfaces = append(interfaces, current)
			}
			current = Interface{Name: strings.TrimPrefix(line, "Interface ")}
		}
		if strings.Contains(line, "type") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				current.Type = parts[1]
			}
		}
	}

	if current.Name != "" {
		interfaces = append(interfaces, current)
	}

	return interfaces, nil
}

// CanRunMonitorMode checks if the system can put an interface into monitor mode
func CanRunMonitorMode() bool {
	// Simple heuristic - in reality this needs more checks
	_, err := exec.LookPath("airmon-ng")
	return err == nil
}
