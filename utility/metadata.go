package utility

import "fmt"

func ConvertToStringMetadata(metadata map[string]interface{}) map[string]string {
	var convertMap = make(map[string]string)
	for key, value := range metadata {
		convertMap[key] = fmt.Sprintf("%v", value)
	}
	return convertMap
}
