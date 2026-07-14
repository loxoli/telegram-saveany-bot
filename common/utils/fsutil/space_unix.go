//go:build !windows

package fsutil

import "golang.org/x/sys/unix"

// AvailableSpace 回傳 path 所在檔案系統目前可用的位元組數（以非特權使用者可用量計）。
func AvailableSpace(path string) (uint64, error) {
	var stat unix.Statfs_t
	if err := unix.Statfs(path, &stat); err != nil {
		return 0, err
	}
	return stat.Bavail * uint64(stat.Bsize), nil
}
