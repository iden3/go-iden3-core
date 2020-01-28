package utils

import (
	"time"
)

func VerifyTimestamp(timestamp int64, timelimit int) bool {
	t := time.Unix(timestamp, 10)
	elapsed := time.Since(t)
	return int(elapsed.Seconds()) <= timelimit
}
