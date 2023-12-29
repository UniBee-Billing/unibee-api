package utility

import "github.com/gogf/gf/v2/encoding/gjson"

func FormatToJson(target interface{}) string {
	if target == nil {
		return ""
	}
	encodeString, err := gjson.EncodeString(target)
	if err != nil {
		return ""
	}
	return encodeString
}
