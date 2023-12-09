package utility

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func GetLineSeparator() string {
	switch runtime.GOOS {
	case "windows":
		return "\r\n"
	default:
		return "\n"
	}
}

func CurrentTimeMillis() (s int64) {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GenerateRandomAlphanumeric(length int) string {
	//rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func JodaTimePrefix() (prefix string) {
	return time.Now().Format("20060102")
}

func CreateMerchantOrderNo() string {
	return fmt.Sprintf("mon%s%s", JodaTimePrefix(), GenerateRandomAlphanumeric(15))
}

func CreateOutRefundNo() string {
	return fmt.Sprintf("orn%s%s", JodaTimePrefix(), GenerateRandomAlphanumeric(15))
}

func GenerateNextInt() int64 {
	//todo mark 工作机器 ID
	return NewSnowflake(1).GenerateID()
}
