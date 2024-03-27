package utility

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

//var (
//	LocalIP       = ""
//	WorkId  int64 = 1
//)
//
//func workIdFromIP() {
//	if WorkId == 1 {
//		LocalIP = DetectLocalIP()
//		WorkIdStr := strings.Replace(LocalIP, ".", "", -1)
//		one, err := strconv.Atoi(WorkIdStr)
//		if err == nil {
//			WorkId = int64(one)
//		}
//	}
//}

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

func CreateSessionId(userId string) string {
	return fmt.Sprintf("us%s%s%s", userId, JodaTimePrefix(), GenerateRandomAlphanumeric(15))
}

func CreateRequestId() string {
	return fmt.Sprintf("req%s%s", JodaTimePrefix(), GenerateRandomAlphanumeric(15))
}

func CreateSubscriptionId() string {
	return fmt.Sprintf("sub%s%s", JodaTimePrefix(), GenerateRandomAlphanumeric(15))
}

func CreateInvoiceId() string {
	return fmt.Sprintf("iv%s%s", JodaTimePrefix(), GenerateRandomAlphanumeric(15))
}

func CreatePendingUpdateId() string {
	return fmt.Sprintf("subup%s%s", JodaTimePrefix(), GenerateRandomAlphanumeric(15))
}

func CreatePaymentId() string {
	return fmt.Sprintf("pay%s%s", JodaTimePrefix(), GenerateRandomAlphanumeric(15))
}

func CreateRefundId() string {
	return fmt.Sprintf("ref%s%s", JodaTimePrefix(), GenerateRandomAlphanumeric(15))
}

// todo mark 高并发情况下生成结果不稳定
//func GenerateNextInt() int64 {
//	workIdFromIP()
//	//return NewSnowflake(WorkId).GenerateID()
//	node, err := snowflake.NewNode(WorkId % 1000)
//	if err != nil {
//		return NewSnowflake(WorkId).GenerateID()
//	} else {
//		return node.Generate().Int64()
//	}
//}

const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateRandomCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateRandomOpenApiKey(length int) (string, error) {
	// Create a byte slice to hold the random bytes
	key := make([]byte, length)

	// Read random bytes from crypto/rand into the byte slice
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	// Encode the random bytes to base64 to get a string representation
	encodedKey := base64.URLEncoding.EncodeToString(key)

	// Truncate the encoded string to the desired length
	// (base64 encoding increases the length by approximately 33%)
	return encodedKey[:length], nil
}
