package utility

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
)

func ConvertToStringMetadata(metadata map[string]interface{}) map[string]string {
	var convertMap = make(map[string]string)
	for key, value := range metadata {
		convertMap[key] = fmt.Sprintf("%v", value)
	}
	return convertMap
}

func MergeMetadata(source string, target map[string]interface{}) map[string]interface{} {
	var metadata = make(map[string]interface{})
	if len(source) > 0 {
		_ = gjson.Unmarshal([]byte(source), &metadata)
	}
	for k, v := range target {
		metadata[k] = v
	}
	return metadata
}
