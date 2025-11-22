package ports

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
)

// ProcessManager handles process operations
type ProcessManager struct{}

// NewProcessManager creates a new process manager
func NewProcessManager() *ProcessManager {
	return &ProcessManager{}
}

// KillProcess kills a process by PID
func (pm *ProcessManager) KillProcess(pid int, force bool) error {
	var cmd *exec.Cmd
	
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		if force {
			cmd = exec.Command("kill", "-9", strconv.Itoa(pid))
		} else {
			cmd = exec.Command("kill", strconv.Itoa(pid))
		}
	} else {
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	
	return cmd.Run()
}

// KillProcessByPort kills a process using a specific port
func (pm *ProcessManager) KillProcessByPort(port int, force bool) error {
	scanner := NewScanner()
	ports, err := scanner.ScanPorts()
	if err != nil {
		return fmt.Errorf("failed to scan ports: %v", err)
	}
	
	for _, p := range ports {
		if p.Port == port {
			return pm.KillProcess(p.PID, force)
		}
	}
	
	return fmt.Errorf("no process found on port %d", port)
}

// GetProcessByPort returns process information for a specific port
func (pm *ProcessManager) GetProcessByPort(port int) (*PortInfo, error) {
	scanner := NewScanner()
	ports, err := scanner.ScanPorts()
	if err != nil {
		return nil, fmt.Errorf("failed to scan ports: %v", err)
	}
	
	for _, p := range ports {
		if p.Port == port {
			return &p, nil
		}
	}
	
	return nil, fmt.Errorf("no process found on port %d", port)
}

