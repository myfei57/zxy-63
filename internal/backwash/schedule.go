package backwash

import "time"

func nowUnix() int64 {
	return time.Now().Unix()
}
