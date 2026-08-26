package osinfo

import (
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var (
	cachedPkgCount string
	pkgCountOnce   sync.Once
)

type Manifest struct {
	Elements []struct {
		PackageAttr *string `json:"packageAttr"`
	} `json:"elements"`
}

func countFromManifest(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return 0
	}
	return len(manifest.Elements)
}

func checkNixSystem() (string, bool) {
	if count := countFromManifest("/nix/var/nix/profiles/system/manifest.json"); count > 0 {
		return strconv.Itoa(count), true
	}
	if files, err := os.ReadDir("/run/current-system/sw/bin"); err == nil && len(files) > 0 {
		return strconv.Itoa(len(files)), true
	}
	return "", false
}

func checkNixUser() (string, bool) {
	homeDir, _ := os.UserHomeDir()
	username := os.Getenv("USER")

	manifestPaths := []string{
		homeDir + "/.local/state/nix/profiles/profile/manifest.json",
		homeDir + "/.local/state/nix/profiles/home-manager/manifest.json",
		"/nix/var/nix/profiles/per-user/" + username + "/profile/manifest.json",
		homeDir + "/.nix-profile/manifest.json",
	}

	for _, mp := range manifestPaths {
		if count := countFromManifest(mp); count > 0 {
			return strconv.Itoa(count), true
		}
	}

	binPaths := []string{
		"/etc/profiles/per-user/" + username + "/bin",
		homeDir + "/.nix-profile/bin",
	}
	for _, bp := range binPaths {
		if files, err := os.ReadDir(bp); err == nil && len(files) > 0 {
			return strconv.Itoa(len(files)), true
		}
	}

	return "", false
}

func GetPackageCount() string {
	pkgCountOnce.Do(func() {
		var parts []string
		if count, ok := checkPacman(); ok {
			parts = append(parts, count+" (pacman)")
		} else if count, ok := checkAPK(); ok {
			parts = append(parts, count+" (apk)")
		} else if count, ok := checkDpkg(); ok {
			parts = append(parts, count+" (dpkg)")
		} else if count, ok := checkRPM(); ok {
			parts = append(parts, count+" (rpm)")
		} else if count, ok := checkXBPS(); ok {
			parts = append(parts, count+" (xbps)")
		} else if count, ok := checkGentoo(); ok {
			parts = append(parts, count+" (portage)")
		} else if count, ok := checkEopkg(); ok {
			parts = append(parts, count+" (eopkg)")
		}

		if sysCount, ok := checkNixSystem(); ok {
			parts = append(parts, sysCount+" (nix-system)")
		}
		if userCount, ok := checkNixUser(); ok {
			parts = append(parts, userCount+" (nix-user)")
		}

		if fCount, ok := checkFlatpak(); ok {
			parts = append(parts, fCount+" (flatpak)")
		}

		if len(parts) > 0 {
			cachedPkgCount = strings.Join(parts, ", ")
		} else {
			cachedPkgCount = "Unknown"
		}
	})
	return cachedPkgCount
}

func checkFlatpak() (string, bool) {
	if _, err := os.Stat("/var/lib/flatpak/app"); err == nil {
		files, err := os.ReadDir("/var/lib/flatpak/app")
		if err == nil && len(files) > 0 {
			return strconv.Itoa(len(files)), true
		}
	}
	return "", false
}

func checkPacman() (string, bool) {
	if _, err := os.Stat("/var/lib/pacman/local"); err == nil {
		files, err := os.ReadDir("/var/lib/pacman/local")
		if err == nil {
			return strconv.Itoa(len(files)), true
		}
	}
	return "", false
}

func checkDpkg() (string, bool) {
	if _, err := os.Stat("/var/lib/dpkg/status"); err == nil {
		out, err := exec.Command("dpkg-query", "-f", ".\n", "-W").Output()
		if err == nil {
			return strconv.Itoa(len(strings.Split(string(out), "\n")) - 1), true
		}
	}
	return "", false
}

func checkGentoo() (string, bool) {
	if _, err := os.Stat("/etc/portage"); err == nil {
		if out, err := exec.Command("eix", "-I", "--only-names").Output(); err == nil {
			outputStr := strings.TrimSpace(string(out))
			if outputStr != "" {
				return strconv.Itoa(len(strings.Split(outputStr, "\n"))), true
			}
		}
		if out, err := exec.Command("qlist", "-I").Output(); err == nil {
			outputStr := strings.TrimSpace(string(out))
			if outputStr != "" {
				return strconv.Itoa(len(strings.Split(outputStr, "\n"))), true
			}
		}
	}
	return "", false
}

func checkRPM() (string, bool) {
	if _, err := os.Stat("/var/lib/rpm"); err == nil {
		out, err := exec.Command("rpm", "-qa").Output()
		if err == nil {
			return strconv.Itoa(len(strings.Split(strings.TrimSpace(string(out)), "\n"))), true
		}
	}
	return "", false
}

func checkXBPS() (string, bool) {
	if _, err := os.Stat("/var/db/xbps"); err == nil {
		files, err := os.ReadDir("/var/db/xbps")
		if err == nil {
			return strconv.Itoa(len(files)), true
		}
	}
	return "", false
}

func checkEopkg() (string, bool) {
	if _, err := os.Stat("/var/lib/eopkg/package"); err == nil {
		files, err := os.ReadDir("/var/lib/eopkg/package")
		if err == nil {
			return strconv.Itoa(len(files)), true
		}
	}
	return "", false
}

func checkAPK() (string, bool) {
	if data, err := os.ReadFile("/lib/apk/db/installed"); err == nil {
		count := 0
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "P:") {
				count++
			}
		}
		if count > 0 {
			return strconv.Itoa(count), true
		}
	}
	return "", false
}
