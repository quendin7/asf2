package hardware

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func GetCPUTemperature() string {
	if runtime.GOOS == "windows" {
		return ""
	}

	files, err := os.ReadDir("/sys/class/hwmon")
	if err != nil {
		return ""
	}

	for _, f := range files {
		namePath := filepath.Join("/sys/class/hwmon", f.Name(), "name")
		name, _ := os.ReadFile(namePath)
		hwmonName := strings.TrimSpace(string(name))

		if hwmonName == "k10temp" || hwmonName == "coretemp" || hwmonName == "zenpower" {
			tempPath := filepath.Join("/sys/class/hwmon", f.Name(), "temp1_input")
			data, err := os.ReadFile(tempPath)
			if err == nil {
				tempRaw := strings.TrimSpace(string(data))
				tempInt, _ := strconv.Atoi(tempRaw)
				return strconv.Itoa(tempInt/1000) + "°C"
			}
		}
	}
	return ""
}
