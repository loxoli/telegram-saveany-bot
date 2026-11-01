package retryutil

import "github.com/duke-git/lancet/v2/retry"

// Retry 與 retry.Retry 相同，但 fn 回傳的錯誤若經 stop 判定為「重試也不會成功」，
// 會立即停止並原樣回傳該錯誤，不再消耗剩餘的重試次數。
func Retry(fn retry.RetryFunc, stop func(error) bool, opts ...retry.Option) error {
	var stopErr error
	err := retry.Retry(func() error {
		err := fn()
		if err != nil && stop(err) {
			stopErr = err
			return nil // 讓 retry.Retry 結束，真正的錯誤由 stopErr 帶出
		}
		return err
	}, opts...)
	if stopErr != nil {
		return stopErr
	}
	return err
}
