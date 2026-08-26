package dodatki

import (
	"os"
	"runtime"
)

func GetUserAndHost() (string, string) {
	var username string
	if runtime.GOOS == "windows" {
		username = os.Getenv("USERNAME")
	} else {
		username = os.Getenv("USER")
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return username, hostname
}
