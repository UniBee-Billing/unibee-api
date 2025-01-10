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

func MergeMetadata(sourceString string, target *map[string]interface{}) map[string]interface{} {
	var metadata = make(map[string]interface{})
	if len(sourceString) > 0 {
		_ = gjson.Unmarshal([]byte(sourceString), &metadata)
	}
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	if target == nil {
		return metadata
	}
	for k, v := range *target {
		metadata[k] = v
	}
	return metadata
}

func MergeStringMetadata(sourceString string, target string) map[string]interface{} {
	var metadata = make(map[string]interface{})
	if len(sourceString) > 0 {
		_ = gjson.Unmarshal([]byte(sourceString), &metadata)
	}
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	if len(target) == 0 {
		return metadata
	}
	var targetMetadata = make(map[string]interface{})
	_ = gjson.Unmarshal([]byte(target), &targetMetadata)
	for k, v := range targetMetadata {
		metadata[k] = v
	}
	return metadata
}
