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
