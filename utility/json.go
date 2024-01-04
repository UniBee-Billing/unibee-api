package utility

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
)

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

func MarshalToJsonString(target interface{}) string {
	if target == nil {
		return ""
	}
	marshal, err := gjson.Marshal(target)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func UnmarshalFromJsonString(target string, one interface{}) error {
	if len(target) > 0 {
		return gjson.Unmarshal([]byte(target), one)
	} else {
		return gerror.New("target is nil")
	}
}
