package hardware

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func GetBatteryInfo() string {
	if runtime.GOOS == "windows" {
		out, err := exec.Command("wmic", "PATH", "Win32_Battery", "GET", "EstimatedChargeRemaining,BatteryStatus", "/format:value").Output()
		if err != nil {
			return "N/A"
		}

		chargeStr := ""
		statusCode := 0
		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "EstimatedChargeRemaining=") {
				chargeStr = strings.TrimSpace(strings.TrimPrefix(line, "EstimatedChargeRemaining="))
			} else if strings.HasPrefix(line, "BatteryStatus=") {
				code, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "BatteryStatus=")))
				if err == nil {
					statusCode = code
				}
			}
		}

		if chargeStr == "" {
			return "N/A"
		}

		statusMap := map[int]string{
			1:  "Discharging",
			2:  "AC",
			3:  "Fully Charged",
			4:  "Low",
			5:  "Critical",
			6:  "Charging",
			7:  "Charging",
			8:  "Charging",
			9:  "Charging",
			10: "Charging",
			11: "Discharging",
		}

		if label, ok := statusMap[statusCode]; ok {
			return fmt.Sprintf("%s%% (%s)", chargeStr, label)
		}
		return fmt.Sprintf("%s%%", chargeStr)
	}

	if runtime.GOOS == "linux" {
		batteryPath := "/sys/class/power_supply/"
		files, err := os.ReadDir(batteryPath)
		if err != nil {
			return "N/A"
		}

		for _, file := range files {
			if strings.HasPrefix(file.Name(), "BAT") {
				capacityFile := filepath.Join(batteryPath, file.Name(), "capacity")
				statusFile := filepath.Join(batteryPath, file.Name(), "status")

				capacity, err := os.ReadFile(capacityFile)
				if err != nil {
					continue
				}
				status, err := os.ReadFile(statusFile)
				if err != nil {
					continue
				}

				capStr := strings.TrimSpace(string(capacity))
				statStr := strings.TrimSpace(string(status))

				return fmt.Sprintf("%s%% (%s)", capStr, statStr)
			}
		}
	}
	return "N/A"
}
