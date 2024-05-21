package util

import (
	"time"
)

func CurrentTimeMillis() int64 {
	return time.Now().UnixNano() / 1e6
}
