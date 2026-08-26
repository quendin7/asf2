package hardware

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	cachedCPUInfo string
	cpuInfoOnce   sync.Once
)

func GetCPUInfo() string {
	cpuInfoOnce.Do(func() {
		cachedCPUInfo = getCPUInfoLinux()
	})
	temp := GetCPUTemperature()
	if temp != "" {
		return cachedCPUInfo + ", " + temp
	}
	return cachedCPUInfo
}

func getCPUInfoLinux() string {
	cpuName := "unknown"
	cpuCores := 0
	cpuThreads := 0
	maxFreqToReport := 0.0

	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "model name") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					cpuName = strings.TrimSpace(parts[1])
					cpuName = regexp.MustCompile(`\s*(\d+-Core\s+)?Processor`).ReplaceAllString(cpuName, "")
					cpuName = strings.Join(strings.Fields(cpuName), " ")
					cpuName = regexp.MustCompile(`@\s*[\d.]+GHz`).ReplaceAllString(cpuName, "")
					if idx := strings.Index(cpuName, " w/"); idx != -1 {
						cpuName = cpuName[:idx]
					}

					cpuName = strings.TrimSpace(cpuName)
				}
			} else if strings.HasPrefix(line, "cpu cores") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					if cores, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
						cpuCores = cores
					}
				}
			} else if strings.HasPrefix(line, "siblings") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					if threads, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
						cpuThreads = threads
					}
				}
			}
		}
	}

	if cpuThreads == 0 {
		cpuThreads = cpuCores
	}
	if cpuThreads == 0 {
		cpuThreads = 1
	}

	for i := 0; i < cpuThreads; i++ {
		freqPath := filepath.Join("/sys/devices/system/cpu", fmt.Sprintf("cpu%d", i), "cpufreq", "cpuinfo_max_freq")
		if data, err := os.ReadFile(freqPath); err == nil {
			if freqKHz, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64); err == nil {
				freqGHz := freqKHz / 1_000_000
				if freqGHz > maxFreqToReport {
					maxFreqToReport = freqGHz
				}
			}
		} else {
			freqPath = filepath.Join("/sys/devices/system/cpu", fmt.Sprintf("cpu%d", i), "cpufreq", "scaling_max_freq")
			if data, err := os.ReadFile(freqPath); err == nil {
				if freqKHz, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64); err == nil {
					freqGHz := freqKHz / 1_000_000
					if freqGHz > maxFreqToReport {
						maxFreqToReport = freqGHz
					}
				}
			}
		}
	}

	var cpuDetails []string
	if cpuName != "unknown" && cpuName != "" {
		cpuDetails = append(cpuDetails, cpuName)
	} else {
		cpuDetails = append(cpuDetails, "Nieznany CPU")
	}

	if cpuCores > 0 {
		coreInfo := fmt.Sprintf("%dC", cpuCores)
		if cpuThreads > 0 && cpuThreads != cpuCores {
			coreInfo += fmt.Sprintf("/%dT", cpuThreads)
		}
		cpuDetails = append(cpuDetails, coreInfo)
	}

	if maxFreqToReport > 0 {
		cpuDetails = append(cpuDetails, fmt.Sprintf("%.2fGHz", maxFreqToReport))
	}

	return strings.Join(cpuDetails, ", ")
}
