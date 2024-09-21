package utility

import "math"

func RoundUp(value float64) int64 {
	return int64(math.Ceil(value))
}
