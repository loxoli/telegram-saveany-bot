package local

import (
	"errors"
	"fmt"
	"io/fs"
	"syscall"
	"testing"
)

func TestIsOutOfSpace(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"預先檢查空間不足", fmt.Errorf("%w: 可用空間 0 B", ErrInsufficientSpace), true},
		{"寫入暫存檔 ENOSPC", fmt.Errorf("failed to copy: %w",
			&fs.PathError{Op: "write", Path: "cache/x", Err: syscall.ENOSPC}), true},
		{"一般錯誤", errors.New("connection reset"), false},
		{"nil", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOutOfSpace(tt.err); got != tt.want {
				t.Errorf("IsOutOfSpace(%v) = %v，預期 %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestFreeSpaceThresholds(t *testing.T) {
	const mib = 1024 * 1024
	orig := availableSpace
	t.Cleanup(func() { availableSpace = orig })

	l := &Local{}
	tests := []struct {
		name     string
		avail    uint64
		beforeOK bool // 下載前檢查（1500 MiB）是否通過
		saveOK   bool // 寫入前檢查（1024 MiB）是否通過
	}{
		{"充足", 2000 * mib, true, true},
		{"剛好 1500 MiB", 1500 * mib, true, true},
		{"1499 MiB：不下載", 1499 * mib, false, true},
		{"剛好 1024 MiB", 1024 * mib, false, true},
		{"1023 MiB：放棄檔案", 1023 * mib, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			availableSpace = func(string) (uint64, error) { return tt.avail, nil }
			if err := l.CheckSpaceBeforeDownload(); (err == nil) != tt.beforeOK {
				t.Errorf("CheckSpaceBeforeDownload() = %v，預期通過 = %v", err, tt.beforeOK)
			} else if err != nil && !IsOutOfSpace(err) {
				t.Errorf("下載前檢查的錯誤應被 IsOutOfSpace 判定為空間不足：%v", err)
			}
			if err := l.checkFreeSpace(minFreeSpace, "寫入前"); (err == nil) != tt.saveOK {
				t.Errorf("寫入前檢查 = %v，預期通過 = %v", err, tt.saveOK)
			}
		})
	}
}
