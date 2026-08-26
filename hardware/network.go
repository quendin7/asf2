package hardware

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

func getWifiSSID(iface string) string {
	cmd := exec.Command("iwconfig", iface)
	output, err := cmd.CombinedOutput()
	if err == nil {
		outputStr := string(output)
		if strings.Contains(outputStr, "ESSID:") {
			parts := strings.Split(outputStr, "ESSID:")
			if len(parts) > 1 {
				ssidParts := strings.Split(parts[1], "\"")
				if len(ssidParts) > 1 {
					return ssidParts[1]
				}
			}
		}
	}
	return ""
}

func GetNetworkInfo() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "unknown"
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		if strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "vboxnet") || strings.HasPrefix(iface.Name, "br-") {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil || len(addrs) == 0 {
			continue
		}

		var ipAddr string
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ipAddr = ipnet.IP.String()
					break
				}
			}
		}

		if ipAddr != "" {
			if strings.HasPrefix(iface.Name, "w") {
				ssid := getWifiSSID(iface.Name)
				if ssid != "" {
					return fmt.Sprintf("%s (%s) [%s]", iface.Name, ipAddr, ssid)
				}
			}
			return fmt.Sprintf("%s (%s)", iface.Name, ipAddr)
		}
	}

	return "No active connection"
}
