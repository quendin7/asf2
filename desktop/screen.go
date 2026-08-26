package desktop

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetScreenResolution() string {

	if os.Getenv("WAYLAND_DISPLAY") != "" {
		if _, err := exec.LookPath("wlr-randr"); err == nil {
			out, err := exec.Command("wlr-randr").Output()
			if err == nil {
				lines := strings.Split(string(out), "\n")
				for _, line := range lines {
					if strings.Contains(line, "current") && strings.Contains(line, "x") {
						fields := strings.Fields(line)
						return fields[0]
					}
				}
			}
		}

		files, err := os.ReadDir("/sys/class/drm")
		if err == nil {
			for _, f := range files {
				modesPath := fmt.Sprintf("/sys/class/drm/%s/modes", f.Name())
				if data, err := os.ReadFile(modesPath); err == nil {
					modes := strings.Split(string(data), "\n")
					if len(modes) > 0 && modes[0] != "" {
						return modes[0]
					}
				}
			}
		}
	}

	if _, err := exec.LookPath("xrandr"); err == nil {
		cmd := exec.Command("xrandr", "--current")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "*") {
					fields := strings.Fields(line)
					if len(fields) > 0 {
						return fields[0]
					}
				}
			}
		}
	}

	return "N/A"
}
