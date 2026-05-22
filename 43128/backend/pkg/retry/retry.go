package retry

import (
	"errors"
	"time"
)

func Do(fn func() error, maxAttempts int, baseDelay time.Duration) error {
	var err error
	for i := 0; i < maxAttempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		if i < maxAttempts-1 {
			delay := baseDelay * time.Duration(1<<uint(i))
			time.Sleep(delay)
		}
	}
	return errors.New("max retry attempts exceeded: " + err.Error())
}
