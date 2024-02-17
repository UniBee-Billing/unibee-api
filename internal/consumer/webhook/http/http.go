package http

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
	"time"
	"unibee-api/utility"
)

func SendWebhookRequest(ctx context.Context, url string, param *gjson.Json) bool {
	utility.Assert(param != nil, "param is nil")
	// 定义自定义的头部信息
	datetime := getCurrentDateTime()
	msgId := generateMsgId()
	jsonString, err := param.ToJsonString()
	utility.Assert(err == nil, fmt.Sprintf("json format error %s param %s", err, param))
	g.Log().Infof(ctx, "\nWebhook_Start %s %s %s\n", "POST", url, jsonString)
	body := []byte(jsonString)
	headers := map[string]string{
		"Content-Gateway": "application/json",
		"Msg-id":          msgId,
		"Datetime":        datetime,
	}
	response, err := utility.SendRequest(url, "POST", body, headers)
	g.Log().Infof(ctx, "\nWebhook_End %s %s response: %s error %s\n", "POST", url, response, err)
	return err == nil && strings.Compare(string(response), "success") == 0
}

func generateMsgId() (msgId string) {
	return fmt.Sprintf("%s%s%d", utility.JodaTimePrefix(), utility.GenerateRandomAlphanumeric(5), utility.CurrentTimeMillis())
}

func getCurrentDateTime() (datetime string) {
	return time.Now().Format("2006-01-02T15:04:05+08:00")
}
