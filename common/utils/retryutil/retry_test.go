package retryutil

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/retry"
)

var errStop = errors.New("stop")

func isStop(err error) bool { return errors.Is(err, errStop) }

func TestRetry(t *testing.T) {
	tests := []struct {
		name      string
		errs      []error // 第 i 次呼叫回傳 errs[i]，超出範圍回傳 nil
		wantCalls int
		wantErr   error
	}{
		{"首次成功", nil, 1, nil},
		{"一般錯誤後成功", []error{errors.New("x")}, 2, nil},
		{"不可重試錯誤立即停止", []error{fmt.Errorf("wrap: %w", errStop)}, 1, errStop},
		{"一般錯誤後遇到不可重試錯誤", []error{errors.New("x"), errStop}, 2, errStop},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			err := Retry(func() error {
				defer func() { calls++ }()
				if calls < len(tt.errs) {
					return tt.errs[calls]
				}
				return nil
			}, isStop, retry.RetryTimes(3), retry.RetryWithLinearBackoff(time.Millisecond))
			if calls != tt.wantCalls {
				t.Errorf("呼叫次數 = %d，預期 %d", calls, tt.wantCalls)
			}
			if tt.wantErr == nil && err != nil {
				t.Errorf("預期成功，得到錯誤 %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("錯誤 = %v，預期包含 %v", err, tt.wantErr)
			}
		})
	}

	t.Run("一般錯誤用完重試次數", func(t *testing.T) {
		calls := 0
		err := Retry(func() error { calls++; return errors.New("x") }, isStop,
			retry.RetryTimes(3), retry.RetryWithLinearBackoff(time.Millisecond))
		if calls != 3 || err == nil {
			t.Errorf("呼叫次數 = %d、錯誤 = %v，預期 3 次且回傳錯誤", calls, err)
		}
	})
}
