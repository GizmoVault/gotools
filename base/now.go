package base

import "time"

type FNNow func() time.Time

func GetNow(now FNNow) time.Time {
	if now == nil {
		return time.Now()
	}

	return now()
}
