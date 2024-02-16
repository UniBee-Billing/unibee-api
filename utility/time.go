package utility

import "time"

func FormatUnixTime(unixTime int64) string {
	// Convert Unix time to time.Time
	timeValue := time.Unix(unixTime, 0)

	// Format time using a layout
	return timeValue.Format("2006-01-02 15:04:05 MST")
}

func MaxInt64(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func MinInt64(a int64, b int64) int64 {
	if a > b {
		return b
	}
	return a
}
