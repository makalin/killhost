package ports

// PortInfo represents information about a port and its process
type PortInfo struct {
	Port        int
	ProcessName string
	PID         int
	Protocol    string
}

