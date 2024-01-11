package utility

import "time"

func FormatUnixTime(unixTime int64) string {
	// Convert Unix time to time.Time
	timeValue := time.Unix(unixTime, 0)

	// Format time using a layout
	return timeValue.Format("2006-01-02 15:04:05 MST")
}
