package dodatki

import (
	"os"
	"runtime"
	"strings"
)

func GetShell() string {
	if runtime.GOOS == "windows" {
		// Try PowerShell first
		if psPath := os.Getenv("PSModulePath"); psPath != "" {
			// PowerShell is available — return "powershell"
			return "powershell"
		}
		if comSpec := os.Getenv("ComSpec"); comSpec != "" {
			parts := strings.Split(comSpec, "\\")
			return parts[len(parts)-1]
		}
		return "unknown"
	}

	shell := os.Getenv("SHELL")
	if shell != "" {
		parts := strings.Split(shell, "/")
		return parts[len(parts)-1]
	}
	return "unknown"
}
