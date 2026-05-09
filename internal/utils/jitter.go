package utils

import (
	"math/rand"
	"time"
)

func Jitter(min, max time.Duration) {
	if min < 0 {
		min = 0
	}
	if max < 0 {
		max = 0
	}
	if min > max {
		min, max = max, min
	}

	delta := max - min
	if delta <= 0 {
		return
	}

	d := min + time.Duration(rand.Int63n(int64(delta)))
	time.Sleep(d)
}
