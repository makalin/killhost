package ports

import (
	"bufio"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// Scanner handles port scanning across different platforms
type Scanner struct{}

// NewScanner creates a new port scanner
func NewScanner() *Scanner {
	return &Scanner{}
}

// ScanPorts scans all active localhost ports and returns port information
func (s *Scanner) ScanPorts() ([]PortInfo, error) {
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		return s.scanUnix()
	}
	return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
}

// scanUnix uses lsof to scan ports on macOS/Linux
func (s *Scanner) scanUnix() ([]PortInfo, error) {
	// Try lsof first (more reliable)
	cmd := exec.Command("lsof", "-i", "-P", "-n", "-sTCP:LISTEN")
	output, err := cmd.Output()
	if err != nil {
		// Fallback to netstat if lsof fails
		return s.scanNetstat()
	}

	return s.parseLsofOutput(string(output))
}

// scanNetstat uses netstat as a fallback
func (s *Scanner) scanNetstat() ([]PortInfo, error) {
	cmd := exec.Command("netstat", "-anv")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to scan ports: %v", err)
	}

	return s.parseNetstatOutput(string(output))
}

// parseLsofOutput parses lsof output
func (s *Scanner) parseLsofOutput(output string) ([]PortInfo, error) {
	var ports []PortInfo
	scanner := bufio.NewScanner(strings.NewReader(output))
	
	// Skip header line
	firstLine := true
	for scanner.Scan() {
		line := scanner.Text()
		if firstLine {
			firstLine = false
			continue
		}
		
		fields := strings.Fields(line)
		if len(fields) < 8 {
			continue
		}

		// lsof format: COMMAND PID USER FD TYPE DEVICE SIZE/OFF NODE NAME
		// Example: node 4473 makalin 12u IPv6 0x... TCP *:3000 (LISTEN)
		// The NAME field can be split: "TCP", "*:3000", "(LISTEN)"
		
		processName := fields[0]
		pidStr := fields[1]
		
		// Find the field containing ":" (the address:port part)
		var addrPortField string
		for _, field := range fields {
			if strings.Contains(field, ":") {
				addrPortField = field
				break
			}
		}
		
		if addrPortField == "" {
			continue
		}
		
		// Extract port from address:port (e.g., *:3000, 127.0.0.1:3000, localhost:3000)
		parts := strings.Split(addrPortField, ":")
		if len(parts) < 2 {
			continue
		}
		
		portStr := strings.TrimSpace(parts[len(parts)-1])
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}
		
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}
		
		// Only include valid ports
		if port > 0 && port < 65536 {
			ports = append(ports, PortInfo{
				Port:        port,
				ProcessName: processName,
				PID:         pid,
				Protocol:    "TCP",
			})
		}
	}
	
	return ports, scanner.Err()
}

// parseNetstatOutput parses netstat output (fallback)
func (s *Scanner) parseNetstatOutput(output string) ([]PortInfo, error) {
	var ports []PortInfo
	scanner := bufio.NewScanner(strings.NewReader(output))
	
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "LISTEN") {
			continue
		}
		
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		
		// netstat format varies by OS
		// macOS: Proto Recv-Q Send-Q Local Address Foreign Address (state) PID/Program
		// Linux: Proto Recv-Q Send-Q Local Address Foreign Address State PID/Program
		
		var localAddr string
		var pid int
		var processName string
		
		// Try to find local address (usually around index 3-4)
		for i, field := range fields {
			if strings.Contains(field, ":") && (strings.Contains(field, "127.0.0.1") || strings.Contains(field, "localhost") || strings.HasPrefix(field, "*:") || strings.HasPrefix(field, ":::")) {
				localAddr = field
				// PID/Program might be in the next few fields
				if i+1 < len(fields) {
					pidProgram := fields[i+1]
					if strings.Contains(pidProgram, "/") {
						parts := strings.Split(pidProgram, "/")
						if len(parts) >= 2 {
							if p, err := strconv.Atoi(parts[0]); err == nil {
								pid = p
								processName = parts[1]
							}
						}
					}
				}
				break
			}
		}
		
		if localAddr == "" {
			// Fallback: try field 3
			if len(fields) > 3 {
				localAddr = fields[3]
			}
		}
		
		if !strings.Contains(localAddr, ":") {
			continue
		}
		
		// Extract port
		parts := strings.Split(localAddr, ":")
		if len(parts) < 2 {
			continue
		}
		
		portStr := parts[len(parts)-1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}
		
		// If we didn't get PID from netstat, try lsof for this specific port
		if pid == 0 {
			pid, processName = s.getProcessInfoByPort(port)
		}
		
		if port > 0 && port < 65536 {
			ports = append(ports, PortInfo{
				Port:        port,
				ProcessName: processName,
				PID:         pid,
				Protocol:    "TCP",
			})
		}
	}
	
	return ports, scanner.Err()
}

// getProcessInfoByPort uses lsof to get PID and process name for a specific port
func (s *Scanner) getProcessInfoByPort(port int) (int, string) {
	cmd := exec.Command("lsof", "-ti", fmt.Sprintf(":%d", port))
	pidOutput, err := cmd.Output()
	if err != nil {
		return 0, "unknown"
	}
	
	pidStr := strings.TrimSpace(string(pidOutput))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, "unknown"
	}
	
	// Get process name
	cmd = exec.Command("ps", "-p", pidStr, "-o", "comm=")
	nameOutput, err := cmd.Output()
	if err != nil {
		return pid, "unknown"
	}
	
	processName := strings.TrimSpace(string(nameOutput))
	if processName == "" {
		processName = "unknown"
	}
	
	return pid, processName
}

