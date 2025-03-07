package utility

import (
	"strconv"
	"strings"
	"unibee/utility/unibee"
)

func ArrayPointJoinToStringPoint(array *[]string) *string {
	if array == nil {
		return nil
	} else {
		return unibee.String(strings.Join(*array, "|"))
	}
}

func JoinToStringPoint(array []string) *string {
	if array == nil {
		return nil
	} else {
		return unibee.String(strings.Join(array, "|"))
	}
}

func SplitToArray(source string) []string {
	if source == "" {
		return make([]string, 0)
	} else {
		return strings.Split(source, "|")
	}
}

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

func IsStringInArray(arr []string, target string) bool {
	if arr == nil || len(arr) == 0 || len(target) == 0 {
		return false
	}
	for _, s := range arr {
		if s == target {
			return true
		}
	}
	return false
}

func IsInt64InArray(arr []int64, target int64) bool {
	if arr == nil || len(arr) == 0 || target == 0 {
		return false
	}
	for _, s := range arr {
		if s == target {
			return true
		}
	}
	return false
}

func IsUint64InArray(arr []uint64, target uint64) bool {
	if arr == nil || len(arr) == 0 || target == 0 {
		return false
	}
	for _, s := range arr {
		if s == target {
			return true
		}
	}
	return false
}

func IsIntInArray(arr []int, target int) bool {
	if arr == nil || len(arr) == 0 || target == 0 {
		return false
	}
	for _, s := range arr {
		if s == target {
			return true
		}
	}
	return false
}
