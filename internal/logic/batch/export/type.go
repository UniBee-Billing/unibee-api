package export

import (
	"context"
	"fmt"
	"time"
	"unibee/utility"
)

func JsonArrayTypeConvert(ctx context.Context, source []interface{}) []int {
	intSlice := make([]int, len(source))
	for i, v := range source {
		if val, ok := v.(float64); ok {
			intSlice[i] = int(val)
		} else {
			utility.Assert(false, fmt.Sprintf("ArrayTypeConvertError from:%v to []int", source))
		}
	}
	return intSlice
}

func JsonArrayTypeConvertUint64(ctx context.Context, source []interface{}) []uint64 {
	intSlice := make([]uint64, len(source))
	for i, v := range source {
		if val, ok := v.(float64); ok {
			intSlice[i] = uint64(val)
		} else {
			utility.Assert(false, fmt.Sprintf("ArrayTypeConvertError from:%v to []int", source))
		}
	}
	return intSlice
}

func GetUTCOffsetFromTimeZone(timeZone string) (int64, error) {
	if len(timeZone) == 0 {
		return 0, fmt.Errorf("time zone is empty")
	}
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		return 0, fmt.Errorf("invalid time zone: %v", err)
	}
	_, offset := time.Now().In(location).Zone()
	return int64(offset), nil
}
