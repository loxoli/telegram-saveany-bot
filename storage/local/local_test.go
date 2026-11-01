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
