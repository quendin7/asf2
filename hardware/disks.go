//go:build linux

package hardware

import (
	"fmt"
	"syscall"
)

func GetDiskInfo() string {
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		return "unknown"
	}

	total := float64(stat.Blocks*uint64(stat.Bsize)) / (1024 * 1024 * 1024)
	free := float64(stat.Bfree*uint64(stat.Bsize)) / (1024 * 1024 * 1024)

	return fmt.Sprintf("%.1f GB / %.1f GB", total-free, total)
}
