package easy

import (
	"time"
)

// Wall-timed execution.
func Timed(action func()) time.Duration {
	start := time.Now()
	action()
	end := time.Now()
	return end.Sub(start)
}
