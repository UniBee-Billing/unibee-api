package utility

import "github.com/gogf/gf/v2/encoding/gjson"

func FormatToJsonString(target interface{}) string {
	if target == nil {
		return ""
	}
	encodeString, err := gjson.EncodeString(target)
	if err != nil {
		return ""
	}
	return encodeString
}

func FormatToGJson(target interface{}) *gjson.Json {
	if target == nil {
		return nil
	}
	jsonData, err := gjson.DecodeToJson(target)
	if err != nil {
		return nil
	}
	return jsonData
}
