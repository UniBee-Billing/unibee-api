package utility

import (
	"strconv"
	"strings"
)

func MergeInt64Arrays(arr1, arr2 []int64) []int64 {
	seen := make(map[int64]bool)
	var result []int64
	for _, num := range arr1 {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	for _, num := range arr2 {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}

	return result
}

func RemoveInt64Arrays(arr, toRemove []int64) []int64 {
	removeMap := make(map[int64]bool)
	for _, num := range toRemove {
		removeMap[num] = true
	}
	var result []int64
	for _, num := range arr {
		if !removeMap[num] {
			result = append(result, num)
		}
	}

	return result
}

func IntListToString(arr []int64) string {
	strArr := make([]string, len(arr))
	for i, num := range arr {
		strArr[i] = strconv.FormatInt(num, 10)
	}
	return strings.Join(strArr, ",")
}
