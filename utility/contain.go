package utility

import "strings"

func IntContainsElement(arr []int, element int) bool {
	for _, val := range arr {
		if val == element {
			return true
		}
	}
	return false
}

func StringContainsElement(arr []string, element string) bool {
	for _, val := range arr {
		if strings.Compare(val, element) == 0 {
			return true
		}
	}
	return false
}
