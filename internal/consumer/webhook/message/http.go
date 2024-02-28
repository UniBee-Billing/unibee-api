package message

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"time"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

func SendWebhookRequest(ctx context.Context, webhookMessage *WebhookMessage) bool {
	utility.Assert(webhookMessage.Data != nil, "param is nil")
	datetime := getCurrentDateTime()
	msgId := generateMsgId()
	jsonString, err := webhookMessage.Data.ToJsonString()
	utility.Assert(err == nil, fmt.Sprintf("json format error %s param %s", err, webhookMessage.Data))
	g.Log().Infof(ctx, "Webhook_Start %s %s %s\n", "POST", webhookMessage.Url, jsonString)
	body := []byte(jsonString)
	headers := map[string]string{
		"Content-Gateway": "application/json",
		"Msg-id":          msgId,
		"Datetime":        datetime,
	}
	response, err := utility.SendRequest(webhookMessage.Url, "POST", body, headers)
	if err != nil {
		g.Log().Errorf(ctx, "Webhook_End %s %s response: %s error %s\n", "POST", webhookMessage.Url, response, err.Error())
	} else {
		g.Log().Infof(ctx, "Webhook_End %s %s response: %s \n", "POST", webhookMessage.Url, response)
	}

	one := &entity.MerchantWebhookLog{
		MerchantId:   webhookMessage.MerchantId,
		EndpointId:   int64(webhookMessage.EndpointId),
		WebhookUrl:   webhookMessage.Url,
		WebhookEvent: string(webhookMessage.Event),
		RequestId:    msgId,
		Body:         jsonString,
		Response:     string(response),
		Mamo:         utility.MarshalToJsonString(err),
		CreateTime:   gtime.Now().Timestamp(),
	}
	_, saveErr := dao.MerchantWebhookLog.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if saveErr != nil {
		g.Log().Errorf(ctx, "Webhook_SaveLog error %s\n", saveErr.Error())
	}

	return err == nil && strings.Compare(string(response), "success") == 0
}

func generateMsgId() (msgId string) {
	return fmt.Sprintf("%s%s%d", utility.JodaTimePrefix(), utility.GenerateRandomAlphanumeric(5), utility.CurrentTimeMillis())
}

func getCurrentDateTime() (datetime string) {
	return time.Now().Format("2006-01-02T15:04:05+08:00")
}
