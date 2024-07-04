package export

import (
	"context"
	"fmt"
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
