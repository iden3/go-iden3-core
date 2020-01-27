package utils

import (
	"time"
)

func VerifyTimestamp(timestamp int64, timelimit int) bool {
	t := time.Unix(timestamp, 10)
	elapsed := time.Since(t)
	if int(elapsed.Seconds()) > timelimit {
		return false
	}
	return true
}
