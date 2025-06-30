package utils

import (
	"time"
)

// NowTimeMilliSecond 获取当前的时间戳（毫秒）
func NowTimeMilliSecond() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}
