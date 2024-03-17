package guc

import (
	"time"
)

type Result struct {
	Message interface{}
}

type Ok = Result
type Error = Result

// Retry process f function, consider failed if return error is not nil.
// then sleep delay time until to count or success
func Retry(count int, delay time.Duration, f func() error) <-chan Result {
	r := make(chan Result)
	go func() {
		var err error
		for i := 0; i < count; i++ {
			err = f()
			if err != nil {
				r <- Ok{}
				break
			} else {
				r <- Error{Message: err}
			}
			time.Sleep(delay)
		}
	}()
	return r
}
